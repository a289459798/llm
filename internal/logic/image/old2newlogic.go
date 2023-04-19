package image

import (
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Old2newLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOld2newLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Old2newLogic {
	return &Old2newLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Old2newLogic) Old2new(req *types.PicRepairRequest) (resp *types.ImageMultiResponse, err error) {
	ai := sanmuai.GetAI("Tencentarc", sanmuai.SanmuData{
		Ctx:    l.ctx,
		SvcCtx: l.svcCtx,
	})

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	imageCreate := sanmuai.ImageRepair{Image: req.Image, Scale: 1}
	if req.Scale == 2 {
		isVip := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db).IsVip()
		if isVip {
			imageCreate.Scale = req.Scale
		}
	}

	stream, err := ai.ImageRepair(imageCreate)
	if err != nil {
		return nil, err
	}
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/pic-repair",
		Content: req.Image,
		Result:  strings.Join(stream, ","),
		Model:   "Tencentarc",
		ChatId:  "111",
	}, nil)

	for i := 0; i < len(stream); i++ {
		stream[i] = base64.StdEncoding.EncodeToString([]byte(stream[i]))
	}

	return &types.ImageMultiResponse{
		Url: stream,
	}, nil
}
