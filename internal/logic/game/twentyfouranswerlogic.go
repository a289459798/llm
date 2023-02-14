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

type TwentyFourAnswerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTwentyFourAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TwentyFourAnswerLogic {
	return &TwentyFourAnswerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TwentyFourAnswerLogic) TwentyFourAnswer(req *types.TwentyFourRequest) (resp *types.GameResponse, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           fmt.Sprintf("这些数字%s玩24点计算小游戏", req.Content),
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
