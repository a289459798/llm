package game

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"context"
	"fmt"

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
	prompt := fmt.Sprintf("这些数字%s玩24点计算小游戏", req.Content)

	stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateCompletion(prompt)
	if err != nil {
		return nil, err
	}
	return &types.GameResponse{Data: utils.TrimHtml(stream.Choices[0].Text)}, nil
}
