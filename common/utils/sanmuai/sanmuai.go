package sanmuai

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
	"net/http"
)

type SanmuAI interface {
	CreateImage(req ImageCreate) (stream []string, err error)
	CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error)
	ImageRepair(image ImageRepair) (result []string, err error)
	ImageRepairAsync(image ImageRepair) (result ImageAsyncTask, err error)
	ImageTask(task ImageAsyncTask) (result ImageTask, err error)
	ImageText(image Image2Text) (result string, err error)
	ImagePS(image ImagePS) (result []string, err error)
	ImagePSAsync(image ImagePS) (result ImageAsyncTask, err error)
}

type SanmuData struct {
	Ctx    context.Context
	SvcCtx *svc.ServiceContext
}

func GetAI(model string, data SanmuData) SanmuAI {
	if model == "Midjourney" || model == "StableDiffusion" {
		return NewJourney(data.Ctx, data.SvcCtx)
	} else if model == "GPT-4" {
		return NewGpt4(data.Ctx, data.SvcCtx)
	} else if model == "Tencentarc" {
		return NewTencentarc(data.Ctx, data.SvcCtx)
	} else if model == "Salesforce" {
		return NewSalesforce(data.Ctx, data.SvcCtx)
	} else if model == "Paintbytext" {
		return NewPaintbytext(data.Ctx, data.SvcCtx)
	} else {
		return NewOpenAi(data.Ctx, data.SvcCtx)
	}
}

func GetProxyIp() string {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://api.xiaoxiangdaili.com/ip/get?appKey=958271663739654144&appSecret=WWYrTxH9&cnt=&wt=json", nil)
	if err != nil {
		return ""
	}
	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	respData := struct {
		Code int `json:"code"`
		Data []struct {
			IP   string `json:"ip"`
			Port int    `json:"port"`
		}
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return ""
	}

	if respData.Code != 200 {
		return ""
	}

	return fmt.Sprintf("http://%s:%d", respData.Data[0].IP, respData.Data[0].Port)
}

func GetKey(db *gorm.DB, channel string) string {
	apikey := &model.Apikey{}
	db.Where("channel = ?", channel).Where("status = ?", 1).Order("rand()").First(&apikey)
	return apikey.Key
}
