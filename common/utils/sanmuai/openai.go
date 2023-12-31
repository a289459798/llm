package sanmuai

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/model"
	"context"
	gogpt "github.com/sashabaranov/go-openai"
	"os"
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

func (ai *OpenAi) getClient(isImage bool) *gogpt.Client {
	apikey := &model.Apikey{}
	tx := ai.SvcCtx.Db.Where("channel = ?", "openai").Where("status = ?", 1)
	if isImage {
		tx.Where("type = ?", 1)
	} else {
		tx.Where("type = ?", 0)
	}
	tx.Order("rand()").Limit(1).Find(apikey)
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

	return ai.getClient(false).CreateCompletionStream(ai.Ctx, gptReq)
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

	return ai.getClient(false).CreateCompletion(ai.Ctx, gptReq)
}

func (ai *OpenAi) CreateChatCompletionStream(content []gogpt.ChatCompletionMessage) (stream *gogpt.ChatCompletionStream, err error) {
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

	return ai.getClient(false).CreateChatCompletionStream(ai.Ctx, gptReq)
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

	return ai.getClient(false).CreateChatCompletion(ai.Ctx, gptReq)
}

func (ai *OpenAi) CreateImage(req ImageCreate) (stream []string, err error) {
	res, err := ai.getClient(true).CreateImage(ai.Ctx, gogpt.ImageRequest{
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

func (ai *OpenAi) CreateEditImage(file *os.File, content string) (stream gogpt.ImageResponse, err error) {
	gptReq := gogpt.ImageEditRequest{
		Image:  file,
		Prompt: content,
		N:      1,
		Size:   "256x256",
	}

	return ai.getClient(true).CreateEditImage(ai.Ctx, gptReq)
}

func (ai *OpenAi) ImageRepair(image ImageRepair) (result []string, err error) {
	return
}

func (ai *OpenAi) ImageText(image Image2Text) (result string, err error) {
	return
}

func (ai *OpenAi) ImageRepairAsync(image ImageRepair) (result ImageAsyncTask, err error) {
	return
}
func (ai *OpenAi) ImageTask(task ImageAsyncTask) (result ImageTask, err error) {
	return
}
func (ai *OpenAi) ImagePS(image ImagePS) (result []string, err error) {
	return
}
func (ai *OpenAi) ImagePSAsync(image ImagePS) (result ImageAsyncTask, err error) {
	return
}
