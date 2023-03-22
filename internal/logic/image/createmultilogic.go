package image

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMultiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMultiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMultiLogic {
	return &CreateMultiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMultiLogic) CreateMulti(req *types.ImageRequest) (resp *types.ImageMultiResponse, err error) {
	valid := utils.Filter(req.Content, l.svcCtx.Db)
	if valid != "" {
		return nil, errors.New(valid)
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	isVip := model.User{ID: uint32(uid)}.Find(l.svcCtx.Db).IsVip()

	prompt := "不能露点，" + req.Content
	imageCreate := sanmuai.ImageCreate{
		Prompt:         prompt,
		N:              1,
		ResponseFormat: "url",
		Size:           "256x256",
	}

	if isVip {
		if req.Model == "gpt-plus" || req.Model == "Midjourney" || req.Model == "StableDiffusion" {
			// 翻译
			message := []gogpt.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "帮我翻译",
				},
				{
					Role:    "user",
					Content: req.Content,
				},
			}
			stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateChatCompletion(message)
			if err == nil && len(stream.Choices) > 0 && stream.Choices[0].Message.Content != "" {
				imageCreate.Prompt = stream.Choices[0].Message.Content
				if req.Model == "Midjourney" {
					imageCreate.Prompt = fmt.Sprintf("midjourney-v4 style %s", stream.Choices[0].Message.Content)
				}
			}
		}
		if req.Number > 0 {
			imageCreate.N = req.Number
		}
		if req.Clarity == "high" {
			imageCreate.Size = "512x512"
		}
	} else {
		req.Model = "gpt"
	}

	ai := sanmuai.GetAI(req.Model, sanmuai.SanmuData{
		Ctx:    l.ctx,
		SvcCtx: l.svcCtx,
	})

	stream, err := ai.CreateImage(imageCreate)
	if err != nil {
		return nil, err
	}

	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/createMulti",
		Content: req.Content,
		Result:  strings.Join(stream, ","),
		Model:   req.Model,
	})

	return &types.ImageMultiResponse{
		Url: stream,
	}, nil
}
