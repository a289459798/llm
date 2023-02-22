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
	prompt := fmt.Sprintf("请随机给我生成一个包含4个数字的算24点的题目，题目只需要有数字，不要包含句号在内的所有字符")

	stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateCompletion(prompt)
	if err != nil {
		return nil, err
	}
	return &types.GameResponse{Data: utils.TrimHtml(stream.Choices[0].Text)}, nil
}
