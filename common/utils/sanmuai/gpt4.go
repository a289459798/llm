package sanmuai

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/model"
	"context"
	gogpt "github.com/sashabaranov/go-openai"
)

type Gpt4 struct {
	Ctx    context.Context
	SvcCtx *svc.ServiceContext
}

func NewGpt4(c context.Context, svcCtx *svc.ServiceContext) *Gpt4 {
	return &Gpt4{
		Ctx:    c,
		SvcCtx: svcCtx,
	}
}

func (ai *Gpt4) getClient() *gogpt.Client {
	apikey := &model.Apikey{}
	ai.SvcCtx.Db.Where("channel = ?", "openai").Where("status = ?", 1).Order("rand()").Limit(1).Find(apikey)
	if apikey.Ori != "" {
		return gogpt.NewOrgClient(apikey.Key, apikey.Ori)
	} else {
		return gogpt.NewClient(apikey.Key)
	}
}

func (ai *Gpt4) CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error) {
	gptReq := gogpt.ChatCompletionRequest{
		Model:            gogpt.GPT4,
		Messages:         content,
		MaxTokens:        1536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}

	return ai.getClient().CreateChatCompletionStream(ai.Ctx, gptReq)
}

func (ai *Gpt4) CreateImage(req ImageCreate) (stream []string, err error) {
	res, err := ai.getClient().CreateImage(ai.Ctx, gogpt.ImageRequest{
		Prompt:         req.Prompt,
		N:              req.N,
		Size:           req.Size,
		ResponseFormat: req.ResponseFormat,
	})
	if err != nil {
		return
	}
	for _, datum := range res.Data {
		stream = append(stream, func() string {
			if req.ResponseFormat == "url" {
				return datum.URL
			}
			return datum.B64JSON
		}())
	}
	return
}
