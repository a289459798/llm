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
	"github.com/zeromicro/go-zero/core/logx"
	"time"
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

	valid := utils.Filter(req.Content, l.svcCtx.Db)
	if valid != "" {
		return nil, errors.New(valid)
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	amount := model.NewAccount(l.svcCtx.Db).GetAccount(uint32(uid), time.Now())
	if amount.ChatAmount-amount.ChatUse < 3 {
		return nil, errors.New("次数已用完")
	}

	prompt := "帮我生成一张图片，图片里面需要包含以下内容：" + req.Content
	stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateImage(sanmuai.ImageCreate{
		Prompt:         prompt,
		N:              1,
		ResponseFormat: "b64_json",
		Size:           "512x512",
	})
	if err != nil {
		return nil, err
	}

	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/create",
		Content: req.Content,
		Result:  "",
	})

	return &types.ImageResponse{
		Url: stream[0],
	}, nil
}
