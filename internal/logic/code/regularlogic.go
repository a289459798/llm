package code

import (
	"chatgpt-tools/common/utils"
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegularLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegularLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegularLogic {
	return &RegularLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegularLogic) Regular(req *types.RegularRequest) (resp *types.CodeResponse, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           fmt.Sprintf("请用以下描述生成一个正则表达式：%s", req.Content),
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
	return &types.CodeResponse{Data: utils.TrimHtml(stream.Choices[0].Text)}, nil
}
