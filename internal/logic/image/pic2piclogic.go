package image

import (
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"os"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Pic2picLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPic2picLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Pic2picLogic {
	return &Pic2picLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Pic2picLogic) Pic2pic(req *types.Pic2picRequest, files map[string][]*multipart.FileHeader) (resp *types.ImageResponse, err error) {
	prompt := req.Style
	if files == nil || len(files["images"]) == 0 {
		return nil, errors.New("请上传图片")
	}
	filename := "data/caches/" + files["images"][0].Filename
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	file, _ := files["images"][0].Open()
	io.Copy(f, file)

	task, err := sanmuai.NewBaiduWX(l.ctx, l.svcCtx).Pic2Pic(&sanmuai.BDImageTaskRequest{
		Prompt: prompt,
		Image:  filename,
	})
	if err != nil {
		return nil, err
	}

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/pic2pic",
		Content: string(task.TaskId),
		Result:  "",
	})

	return &types.ImageResponse{
		Url: string(task.TaskId),
	}, nil
}
