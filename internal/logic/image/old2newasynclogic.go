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

	stream, err := ai.ImageRepairAsync(sanmuai.ImageRepair{Image: req.Image})
	if err != nil {
		return nil, err
	}

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
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
