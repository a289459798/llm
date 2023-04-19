package image

import (
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

type Old2newAsyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOld2newAsyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Old2newAsyncLogic {
	return &Old2newAsyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Old2newAsyncLogic) Old2newAsync(req *types.PicRepairRequest) (resp *types.ImageMultiAsyncResponse, err error) {
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

	stream, err := ai.ImageRepairAsync(imageCreate)
	if err != nil {
		return nil, err
	}

	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/pic-repair",
		Content: req.Image,
		ChatId:  stream.Task,
		Model:   "Tencentarc",
	}, nil)

	return &types.ImageMultiAsyncResponse{
		Model: "Tencentarc",
		Task:  stream.Task,
	}, nil
}
