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
	"strconv"
	"strings"
	"time"
)

type Journey struct {
	Ctx    context.Context
	SvcCtx *svc.ServiceContext
}

func NewJourney(c context.Context, svcCtx *svc.ServiceContext) *Journey {
	return &Journey{
		Ctx:    c,
		SvcCtx: svcCtx,
	}
}

func (ai *Journey) CreateImage(image ImageCreate) (result []string, err error) {
	authorization := GetKey(ai.SvcCtx.Db, "replicate")

	uuid, err := ai.create(authorization, image)
	if err != nil {
		return
	}

	// 获取信息
	resultChan := make(chan []string)
	quitChan := make(chan string)
	timeout := time.After(180 * time.Second)
	timer := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-quitChan:
				return
			case <-timer.C:
				go func(resultChan chan []string, quitChan chan string) {
					client := &http.Client{}
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.replicate.com/v1/predictions/%s", uuid), nil)
					if err != nil {
						return
					}
					req.Header.Set("Authorization", "Token "+authorization)
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
						Uuid      string   `json:"uuid"`
						Error     string   `json:"error"`
						Status    string   `json:"status"`
					}{}
					if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
						return
					}
					fmt.Println(respData)
					if respData.Error != "" {
						resultChan <- respData.Output
						close(quitChan)
						return
					}
					if respData.Output != nil {
						for i := 0; i < len(respData.Output); i++ {
							respData.Output[i] = strings.Replace(respData.Output[i], "https://replicate.delivery/", "https://img2.smuai.com/", 1)
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

	if len(result) == 0 {
		err = errors.New("系统错误请重试")
	}
	return
}

func (ai *Journey) CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error) {
	err = errors.New("该模型不支持会话")
	return
}

func (ai *Journey) ImageRepair(image ImageRepair) (result []string, err error) {
	return
}

func (ai *Journey) create(cookie string, image ImageCreate) (uuid string, err error) {
	client := &http.Client{}
	s := strings.Split(image.Size, "x")
	width, _ := strconv.Atoi(s[0])
	height, _ := strconv.Atoi(s[1])
	if width > 768 {
		width = 768
	}

	if height > 768 {
		height = 768
	}
	data := map[string]interface{}{
		"version": "9936c2001faa2194a261c01381f90e65261879985476014a0a37a334593a05eb",
		"input": map[string]interface{}{
			"guidance_scale":      7,
			"width":               width,
			"height":              height,
			"num_inference_steps": 50,
			"num_outputs":         image.N,
			"prompt":              image.Prompt,
		},
	}

	// 将数据转换为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.replicate.com/v1/predictions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errors.New("错误，请重试1")
	}

	req.Header.Set("Authorization", "Token "+cookie)
	req.Header.Set("Content-Type", "application/json")
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", errors.New(resp.Status)
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

func getCookie() (cookieString string, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://replicate.com/prompthero/openjourney", nil)
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

func (ai *Journey) ImageText(image Image2Text) (result string, err error) {
	return
}

func (ai *Journey) ImageRepairAsync(image ImageRepair) (result ImageAsyncTask, err error) {
	return
}
func (ai *Journey) ImageTask(task ImageAsyncTask) (result ImageTask, err error) {
	return
}

func (ai *Journey) ImagePS(image ImagePS) (result []string, err error) {
	return
}
func (ai *Journey) ImagePSAsync(image ImagePS) (result ImageAsyncTask, err error) {
	return
}
