package image

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"unicode"
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
	isVip := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db).IsVip()

	prompt := req.Content
	imageCreate := sanmuai.ImageCreate{
		Prompt:         prompt,
		N:              1,
		ResponseFormat: "url",
		Size:           "512x512",
	}

	paramsMap := make(map[string]interface{})
	paramsMap["number"] = 1

	hasChinese := false
	for _, r := range req.Content {
		if unicode.Is(unicode.Han, r) {
			hasChinese = true
			break
		}
	}

	if hasChinese {
		// 翻译
		message := []gogpt.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "帮我翻译",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("你化身为翻译官，把我提供内容翻译成英文，注意如果已经是英文了直接回答我原文，不要回复其他内容，你要翻译的内容是：%s", req.Content),
			},
		}
		stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateChatCompletion(message)
		if err != nil {
			return nil, err
		}
		if len(stream.Choices) > 0 && stream.Choices[0].Message.Content != "" {
			imageCreate.Prompt = stream.Choices[0].Message.Content
		}
	} else {
		imageCreate.Prompt = req.Content
	}

	if req.Model == "Midjourney" {
		imageCreate.Prompt = fmt.Sprintf("mdjrny-v4 style %s", imageCreate.Prompt)
	}

	if isVip {
		if req.Number > 0 {
			imageCreate.N = req.Number
			paramsMap["number"] = req.Number
		}
		paramsMap["clarity"] = req.Clarity
		if req.Clarity == "high" {
			imageCreate.Size = "512x512"
		} else if req.Clarity == "superhigh" {
			imageCreate.Size = "1024x1024"
			imageCreate.Prompt += " 8k"
		}
	} else {
		req.Model = "GPT-PLUS"
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
		Uid:         uint32(uid),
		Type:        "image/createMulti",
		ShowContent: req.Content,
		Content:     imageCreate.Prompt,
		Result:      fmt.Sprintf("%s|||%s", imageCreate.Prompt, strings.Join(stream, ",")),
		Model:       req.Model,
		ChatId:      "111",
	}, &service.RecordParams{
		Params: func() string {
			if paramsMap != nil {
				params, _ := json.Marshal(paramsMap)
				return string(params)
			}
			return ""
		}(),
	})

	if req.Model != "" {
		for i := 0; i < len(stream); i++ {
			stream[i] = base64.StdEncoding.EncodeToString([]byte(stream[i]))
		}
	}

	return &types.ImageMultiResponse{
		Url: stream,
	}, nil
}
