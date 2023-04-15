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
	"net/url"
	"strings"
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
	timeout := time.After(300 * time.Second)
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
					if respData.Prediction.Output != nil && len(respData.Prediction.Output) > 0 {
						for i := 0; i < len(respData.Prediction.Output); i++ {
							respData.Prediction.Output[i] = strings.Replace(respData.Prediction.Output[i], "https://replicate.delivery/", "https://img2.smuai.com/", 1)
						}
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
	ip := GetProxyIp()
	client := &http.Client{}
	if ip != "" {
		proxyUrl, _ := url.Parse(ip)
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	data := map[string]map[string]interface{}{
		"inputs": {
			"img":     image.Image,
			"scale":   2,
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
		return "", errors.New("错误，请重试1")
	}
	req.Header.Set("x-csrftoken", cookie)
	req.Header.Set("Content-Type", "application/json")
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("错误，请重试2")
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

func (ai *Tencentarc) ImageText(image Image2Text) (result string, err error) {
	return
}

func (ai *Tencentarc) ImageRepairAsync(image ImageRepair) (result ImageAsyncTask, err error) {
	cookie, err := getTencentarcCookie()
	if cookie == "" {
		err = errors.New("cookie is empty")
		return
	}
	uuid, err := createRepair(cookie, image)
	result.Task = uuid
	result.Session = cookie
	return
}
func (ai *Tencentarc) ImageTask(task ImageAsyncTask) (result ImageTask, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://replicate.com/api/models/tencentarc/gfpgan/versions/9283608cc6b7be6b65a8e44983db012355fde4132009bf99d976b2f0896856a3/predictions/%s", task.Task), nil)
	fmt.Println(fmt.Sprintf("https://replicate.com/api/models/tencentarc/gfpgan/versions/9283608cc6b7be6b65a8e44983db012355fde4132009bf99d976b2f0896856a3/predictions/%s", task.Task))
	if err != nil {
		return
	}
	//req.Header.Set("x-csrftoken", task.Session)
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
	fmt.Println(respData)
	result.Output = respData.Prediction.Output
	return
}

func (ai *Tencentarc) ImagePS(image ImagePS) (result []string, err error) {
	return
}
func (ai *Tencentarc) ImagePSAsync(image ImagePS) (result ImageAsyncTask, err error) {
	return
}
