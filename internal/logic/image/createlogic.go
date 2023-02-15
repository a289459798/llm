package image

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ImageRequest) (resp *types.ImageResponse, err error) {
	gptReq := gogpt.ImageRequest{
		Prompt: "帮我画一幅画，需要包含下一下内容：" + req.Content,
		N:      1,
	}
	ctx := context.Background()
	stream, err := l.svcCtx.GptClient.CreateImage(ctx, gptReq)
	if err != nil {
		return nil, err
	}

	return &types.ImageResponse{
		Url: stream.Data[0].URL,
	}, nil
}
