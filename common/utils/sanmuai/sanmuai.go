package sanmuai

import (
	"chatgpt-tools/internal/svc"
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"
	"math/rand"
	"time"
)

type SanmuAI interface {
	CreateImage(req ImageCreate) (stream []string, err error)
	CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error)
	ImageRepair(image ImageRepair) (result []string, err error)
	ImageRepairAsync(image ImageRepair) (result ImageAsyncTask, err error)
	ImageTask(task ImageAsyncTask) (result ImageTask, err error)
	ImageText(image Image2Text) (result string, err error)
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
	} else {
		return NewOpenAi(data.Ctx, data.SvcCtx)
	}
}

func GetProxyIp() string {
	ips := []string{
		"http://27.159.66.131:11054",
		"http://27.159.66.131:11206",
		"http://27.159.66.131:11091",
		"http://27.159.66.131:11220",
		"http://27.159.66.131:11050",
		"",
	}
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(6)
	fmt.Println(randomNum)
	if randomNum < len(ips) {
		return ips[randomNum]
	}
	return ""
}
