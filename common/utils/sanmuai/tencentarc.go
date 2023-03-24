package sanmuai

import (
	"bytes"
	"chatgpt-tools/internal/svc"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"
	"net/http"
	"time"
)

type Tencentarc struct {
	Ctx    context.Context
	SvcCtx *svc.ServiceContext
}

func NewTencentarc(c context.Context, svcCtx *svc.ServiceContext) *Tencentarc {
	return &Tencentarc{
		Ctx:    c,
		SvcCtx: svcCtx,
	}
}

func (ai *Tencentarc) ImageRepair(image ImageRepair) (result []string, err error) {
	cookie, err := getTencentarcCookie()
	if cookie == "" {
		err = errors.New("cookie is empty")
		return
	}

	uuid, err := createRepair(cookie, image)
	if err != nil {
		return
	}

	// 获取信息
	resultChan := make(chan []string)
	quitChan := make(chan string)
	timeout := time.After(60 * time.Second)
	timer := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-quitChan:
				return
			case <-timer.C:
				go func(resultChan chan []string, quitChan chan string) {
					client := &http.Client{}
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://replicate.com/api/models/tencentarc/gfpgan/versions/9283608cc6b7be6b65a8e44983db012355fde4132009bf99d976b2f0896856a3/predictions/%s", uuid), nil)
					if err != nil {
						return
					}
					req.Header.Set("x-csrftoken", cookie)
					req.Header.Set("Content-Type", "application/json")
					// 发送请求并获取响应
					resp, err := client.Do(req)
					if err != nil {
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode != 200 {
						err = errors.New(resp.Status)
						return
					}

					respData := struct {
						Prediction struct {
							Output    []string `json:"output_files"`
							CreatedAt string   `json:"created_at"`
							Uuid      string   `json:"uuid"`
							Error     string   `json:"error"`
							Status    string   `json:"status"`
						} `json:"prediction"`
					}{}
					if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
						return
					}
					if respData.Prediction.Output != nil {
						resultChan <- respData.Prediction.Output
						close(quitChan)
					}
				}(resultChan, quitChan)
			}
		}
	}()

	select {
	case result = <-resultChan:
	case <-timeout:
		close(quitChan)
		err = errors.New("timeout")
	}
	timer.Stop()
	return
}

func (ai *Tencentarc) CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error) {
	err = errors.New("该模型不支持会话")
	return
}

func (ai *Tencentarc) CreateImage(image ImageCreate) (result []string, err error) {
	return
}

func createRepair(cookie string, image ImageRepair) (uuid string, err error) {
	client := &http.Client{}
	data := map[string]map[string]interface{}{
		"inputs": {
			"img":     image.Image,
			"scale":   1,
			"version": "v1.4",
		},
	}

	// 将数据转换为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, "https://replicate.com/api/models/tencentarc/gfpgan/versions/9283608cc6b7be6b65a8e44983db012355fde4132009bf99d976b2f0896856a3/predictions", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	req.Header.Set("x-csrftoken", cookie)
	req.Header.Set("Content-Type", "application/json")
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		err = errors.New(resp.Status)
		return
	}

	respData := struct {
		UUID string `json:"uuid"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return
	}
	uuid = respData.UUID
	return
}

func getTencentarcCookie() (cookieString string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://replicate.com/tencentarc/gfpgan", nil)
	if err != nil {
		return
	}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	// 获取所有 cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "csrftoken" {
			cookieString = cookie.Value
			break
		}
	}
	return
}
