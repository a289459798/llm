package image

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WatermarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWatermarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WatermarkLogic {
	return &WatermarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WatermarkLogic) Watermark(req *types.WatermarkRequest) (resp *types.ImageResponse, err error) {
	gptReq := gogpt.ImageEditRequest{
		Image: ,
		Prompt: fmt.Sprintf("在文字%s方加上文字大小为%d颜色为%s透明度为%f的水印，水印内容为:%s", req.Position, req.FontSize, req.Color, req.Opacity, req.Content),
		N:      1,
	}
	ctx := context.Background()
	stream, err := l.svcCtx.GptClient.CreateEditImage(ctx, gptReq)
	if err != nil {
		return nil, err
	}

	return &types.ImageResponse{
		Url: stream.Data[0].URL,
	}, nil
}
