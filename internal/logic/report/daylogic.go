package report

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DayLogic {
	return &DayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DayLogic) Day(req *types.ReportRequest) (resp *types.ReportResponse, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           "请帮我把以下的工作内容填充为一篇完整的日报,用 html 格式以分点叙述的形式输出:" + req.Content,
		MaxTokens:        1536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}
	ctx := context.Background()
	stream, err := l.svcCtx.GptClient.CreateCompletion(ctx, gptReq)
	if err != nil {
		return nil, err
	}

	fmt.Println(stream.Choices)
	return &types.ReportResponse{
		Data: stream.Choices[0].Text,
	}, nil
}
