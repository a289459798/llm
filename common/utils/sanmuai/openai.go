package sanmuai

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/model"
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type OpenAi struct {
	Ctx    context.Context
	SvcCtx *svc.ServiceContext
}

func NewOpenAi(c context.Context, svcCtx *svc.ServiceContext) *OpenAi {
	return &OpenAi{
		Ctx:    c,
		SvcCtx: svcCtx,
	}
}

func (ai *OpenAi) getClient() *gogpt.Client {
	apikey := &model.Apikey{}
	ai.SvcCtx.Db.Where("channel = ?", "openai").Where("status = ?", 1).Order("rand()").Limit(1).Find(apikey)
	if apikey.Ori != "" {
		return gogpt.NewOrgClient(apikey.Key, apikey.Ori)
	} else {
		return gogpt.NewClient(apikey.Key)
	}
}

func (ai *OpenAi) CreateCompletionStream(content string) (stream *gogpt.CompletionStream, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           content,
		MaxTokens:        1536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}

	return ai.getClient().CreateCompletionStream(ai.Ctx, gptReq)
}

func (ai *OpenAi) CreateCompletion(content string) (stream gogpt.CompletionResponse, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           content,
		MaxTokens:        1536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}

	return ai.getClient().CreateCompletion(ai.Ctx, gptReq)
}

func (ai *OpenAi) CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error) {
	gptReq := gogpt.ChatCompletionRequest{
		Model:            gogpt.GPT3Dot5Turbo,
		Messages:         content,
		MaxTokens:        2536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}

	return ai.getClient().CreateChatCompletionStream(ai.Ctx, gptReq)
}

func (ai *OpenAi) CreateChatCompletion(content []gogpt.ChatCompletionMessage) (stream gogpt.ChatCompletionResponse, err error) {
	gptReq := gogpt.ChatCompletionRequest{
		Model:            gogpt.GPT3Dot5Turbo,
		Messages:         content,
		MaxTokens:        1536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}

	return ai.getClient().CreateChatCompletion(ai.Ctx, gptReq)
}

func (ai *OpenAi) CreateImage(content string) (stream gogpt.ImageResponse, err error) {
	gptReq := gogpt.ImageRequest{
		Prompt:         content,
		N:              1,
		ResponseFormat: "b64_json",
		Size:           "256x256",
	}

	return ai.getClient().CreateImage(ai.Ctx, gptReq)
}

func (ai *OpenAi) CreateEditImage(content string) (stream gogpt.ImageResponse, err error) {
	gptReq := gogpt.ImageEditRequest{
		Prompt: content,
		N:      1,
		Size:   "256x256",
	}

	return ai.getClient().CreateEditImage(ai.Ctx, gptReq)
}
