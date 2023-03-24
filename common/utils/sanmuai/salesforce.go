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

type Salesforce struct {
	Ctx    context.Context
	SvcCtx *svc.ServiceContext
}

func NewSalesforce(c context.Context, svcCtx *svc.ServiceContext) *Salesforce {
	return &Salesforce{
		Ctx:    c,
		SvcCtx: svcCtx,
	}
}

func (ai *Salesforce) ImageRepair(image ImageRepair) (result []string, err error) {
	return
}

func (ai *Salesforce) CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error) {
	err = errors.New("该模型不支持会话")
	return
}

func (ai *Salesforce) CreateImage(image ImageCreate) (result []string, err error) {
	return
}

func (ai *Salesforce) ImageText(image Image2Text) (result string, err error) {
	cookie, err := getSalesforceCookie()
	if cookie == "" {
		err = errors.New("cookie is empty")
		return
	}

	uuid, err := createText(cookie, image)
	if err != nil {
		return
	}

	// 获取信息
	resultChan := make(chan string)
	quitChan := make(chan string)
	timeout := time.After(60 * time.Second)
	timer := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-quitChan:
				return
			case <-timer.C:
				go func(resultChan chan string, quitChan chan string) {
					client := &http.Client{}
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://replicate.com/api/models/salesforce/blip/versions/2e1dddc8621f72155f24cf2e0adbde548458d3cab9f00c0139eea840d0ac4746/predictions/%s", uuid), nil)
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
							Output    string `json:"output"`
							CreatedAt string `json:"created_at"`
							Uuid      string `json:"uuid"`
							Error     string `json:"error"`
							Status    string `json:"status"`
						} `json:"prediction"`
					}{}
					if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
						return
					}
					if respData.Prediction.Output != "" {
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

func createText(cookie string, image Image2Text) (uuid string, err error) {
	client := &http.Client{}
	data := map[string]map[string]interface{}{
		"inputs": {
			"image": image.Image,
			"task":  "image_captioning",
		},
	}

	// 将数据转换为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, "https://replicate.com/api/models/salesforce/blip/versions/2e1dddc8621f72155f24cf2e0adbde548458d3cab9f00c0139eea840d0ac4746/predictions", bytes.NewBuffer(jsonData))
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

func getSalesforceCookie() (cookieString string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://replicate.com/salesforce/blip", nil)
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
