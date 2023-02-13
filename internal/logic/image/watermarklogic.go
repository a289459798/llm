package image

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"image"
	"image/png"
	"net/http"
	"os"

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

	imgResp, err := http.Get(req.Url)
	if err != nil {
		return nil, err
	}
	defer imgResp.Body.Close()

	// 解码图片
	img, _, err := image.Decode(imgResp.Body)
	if err != nil {
		return nil, err
	}

	out, err := os.Create("tmp/image.png")
	if err != nil {
		return nil, err
	}
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		return nil, err
	}

	//defer os.Remove("tmp/image.png")

	gptReq := gogpt.ImageEditRequest{
		Image:  out,
		Mask:   out,
		Prompt: fmt.Sprintf("在文字%s方加上文字大小为%d颜色为%s透明度为%f的水印，水印内容为:%s", req.Position, req.FontSize, req.Color, req.Opacity, req.Content),
		N:      1,
		Size:   "512x512",
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
