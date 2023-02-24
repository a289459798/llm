package image

import (
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type Pic2pictaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPic2pictaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Pic2pictaskLogic {
	return &Pic2pictaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Pic2pictaskLogic) Pic2pictask(req *types.Pic2picTaskRequest) (resp *types.ImageResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	task := req.Task
	record := &model.Record{}
	l.svcCtx.Db.Where("uid = ?", uid).Where("type = ?", "image/pic2pic").Where("content = ?", task).Find(&record)
	if record.Content == "" {
		return nil, errors.New("图片任务不存在")
	}

	if record.Result != "" {
		return &types.ImageResponse{
			Url: record.Result,
		}, nil
	}

	img, err := sanmuai.NewBaiduWX(l.ctx, l.svcCtx).Pic2PicTask(task)
	if err != nil {
		return nil, err
	}

	if img.Status == 0 {
		return &types.ImageResponse{
			Url: "",
		}, nil
	}

	record.Result = img.ImgUrls[0]["image"]
	l.svcCtx.Db.Save(record)
	return &types.ImageResponse{
		Url: record.Result,
	}, nil
}
