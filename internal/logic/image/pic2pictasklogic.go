package image

import (
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
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

	task := req.Task
	record := &model.Pic2Pic{}
	l.svcCtx.Db.Where("task_id = ?", task).Find(&record)
	if record.TaskId == "" {
		return nil, errors.New("图片任务不存在")
	}

	if record.Url != "" {
		return &types.ImageResponse{
			Url: record.Url,
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

	record.Url = img.ImgUrls[0]["image"]
	l.svcCtx.Db.Save(record)
	return &types.ImageResponse{
		Url: record.Url,
	}, nil
}
