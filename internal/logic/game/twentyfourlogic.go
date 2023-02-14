package game

import (
	"chatgpt-tools/common/utils"
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TwentyFourLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTwentyFourLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TwentyFourLogic {
	return &TwentyFourLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TwentyFourLogic) TwentyFour() (resp *types.GameResponse, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           fmt.Sprintf("请随机给我生成一个包含4个数字的算24点的题目，题目只需要有数字，不要包含句号在内的所有字符"),
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
