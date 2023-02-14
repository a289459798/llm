package game

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/zeromicro/go-zero/core/logx"
)

type IdiomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdiomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdiomLogic {
	return &IdiomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdiomLogic) Idiom() (resp *types.GameResponse, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           "请随机给我生成一个四字成语，不要包含任何符号",
		MaxTokens:        1536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}

	stream, err := l.svcCtx.GptClient.CreateCompletion(l.ctx, gptReq)
	if err != nil {
		return nil, err
	}
	return &types.GameResponse{Data: utils.TrimHtml(stream.Choices[0].Text)}, nil
}
