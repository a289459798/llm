package image

import (
	"chatgpt-tools/model"
	"chatgpt-tools/service"
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
		Prompt:         "帮我画一幅画，需要包含以下内容：" + req.Content,
		N:              1,
		ResponseFormat: "b64_json",
		Size:           "512x512",
	}
	ctx := context.Background()
	stream, err := l.svcCtx.GptClient.CreateImage(ctx, gptReq)
	if err != nil {
		return nil, err
	}

	service.NewRecord(l.svcCtx.Db).Insert(model.Record{
		Uid:     l.ctx.Value("uid").(uint32),
		Type:    "image/create",
		Content: "",
		Result:  "",
	})

	return &types.ImageResponse{
		Url: stream.Data[0].B64JSON,
	}, nil
}
