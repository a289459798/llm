package sanmuai

import (
	"bytes"
	"chatgpt-tools/internal/svc"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	cookie, err := getCookie()
	if cookie == "" {
		err = errors.New("cookie is empty")
		return
	}

	uuid, err := create(cookie, image)
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
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://replicate.com/api/models/prompthero/openjourney/versions/9936c2001faa2194a261c01381f90e65261879985476014a0a37a334593a05eb/predictions/%s", uuid), nil)
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
							Output    []string `json:"output"`
							CreatedAt string   `json:"created_at"`
						} `json:"prediction"`
					}{}
					if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
						return
					}
					fmt.Println(respData.Prediction.Output)
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

func create(cookie string, image ImageCreate) (uuid string, err error) {
	client := &http.Client{}
	s := strings.Split(image.Size, "x")
	width, _ := strconv.Atoi(s[0])
	height, _ := strconv.Atoi(s[1])
	data := map[string]map[string]interface{}{
		"inputs": {
			"guidance_scale":      7,
			"width":               width,
			"height":              height,
			"num_inference_steps": 50,
			"num_outputs":         image.N,
			"prompt":              image.Prompt,
			"seed":                nil,
		},
	}

	// 将数据转换为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, "https://replicate.com/api/models/prompthero/openjourney/versions/9936c2001faa2194a261c01381f90e65261879985476014a0a37a334593a05eb/predictions", bytes.NewBuffer(jsonData))
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