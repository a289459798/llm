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

type Paintbytext struct {
	Ctx    context.Context
	SvcCtx *svc.ServiceContext
}

func NewPaintbytext(c context.Context, svcCtx *svc.ServiceContext) *Paintbytext {
	return &Paintbytext{
		Ctx:    c,
		SvcCtx: svcCtx,
	}
}

func (ai *Paintbytext) ImagePS(image ImagePS) (result []string, err error) {

	uuid, err := createPS(image)
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
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://paintbytext.chat/api/predictions/%s", uuid), nil)
					if err != nil {
						return
					}
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
						Output    []string `json:"output"`
						CreatedAt string   `json:"created_at"`
						ID        string   `json:"id"`
						Error     string   `json:"error"`
						Status    string   `json:"status"`
					}{}
					if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
						return
					}
					if respData.Output != nil && len(respData.Output) > 0 {
						for i := 0; i < len(respData.Output); i++ {
							respData.Output[i] = strings.Replace(respData.Output[i], "https://replicate.delivery/", "http://img2.smuai.com/", 1)
						}
						resultChan <- respData.Output
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

func (ai *Paintbytext) CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error) {
	err = errors.New("该模型不支持会话")
	return
}

func (ai *Paintbytext) CreateImage(image ImageCreate) (result []string, err error) {
	return
}

func createPS(image ImagePS) (uuid string, err error) {
	ip := GetProxyIp()
	fmt.Println("ip", ip)
	client := &http.Client{}
	if ip != "" {
		proxyUrl, _ := url.Parse(ip)
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	data := map[string]string{
		"image":  image.Image,
		"prompt": image.Text,
	}

	// 将数据转换为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, "https://paintbytext.chat/api/predictions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errors.New("错误，请重试1")
	}
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
		ID string `json:"id"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return
	}
	uuid = respData.ID
	return
}

func (ai *Paintbytext) ImageText(image Image2Text) (result string, err error) {
	return
}

func (ai *Paintbytext) ImagePSAsync(image ImagePS) (result ImageAsyncTask, err error) {
	uuid, err := createPS(image)
	result.Task = uuid
	return
}
func (ai *Paintbytext) ImageTask(task ImageAsyncTask) (result ImageTask, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://paintbytext.chat/api/predictions/%s", task.Task), nil)
	if err != nil {
		return
	}
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
		Output    []string `json:"output_files"`
		CreatedAt string   `json:"created_at"`
		ID        string   `json:"id"`
		Error     string   `json:"error"`
		Status    string   `json:"status"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return
	}
	fmt.Println(respData)
	result.Output = respData.Output
	return
}

func (ai *Paintbytext) ImageRepair(image ImageRepair) (result []string, err error) {
	return
}
func (ai *Paintbytext) ImageRepairAsync(image ImageRepair) (result ImageAsyncTask, err error) {
	return
}
