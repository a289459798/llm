package image

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMultiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMultiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMultiLogic {
	return &CreateMultiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMultiLogic) CreateMulti(req *types.ImageRequest) (resp *types.ImageMultiResponse, err error) {
	valid := utils.Filter(req.Content, l.svcCtx.Db)
	if valid != "" {
		return nil, errors.New(valid)
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	prompt := "帮我生成一张图片，图片里面需要包含以下内容：" + req.Content

	imageCreate := sanmuai.ImageCreate{
		Prompt:         prompt,
		N:              1,
		ResponseFormat: "url",
		Size:           "512x512",
	}
	ai := sanmuai.GetAI(req.Model, sanmuai.SanmuData{
		Ctx:    l.ctx,
		SvcCtx: l.svcCtx,
	})

	stream, err := ai.CreateImage(imageCreate)
	if err != nil {
		return nil, err
	}

	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/createMulti",
		Content: req.Content,
		Result:  strings.Join(stream, ","),
	})

	return &types.ImageMultiResponse{
		Url: stream,
	}, nil
}