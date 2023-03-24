package sanmuai

import (
	"chatgpt-tools/internal/svc"
	"context"
	gogpt "github.com/sashabaranov/go-openai"
)

type SanmuAI interface {
	CreateImage(req ImageCreate) (stream []string, err error)
	CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error)
	ImageRepair(image ImageRepair) (result []string, err error)
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
