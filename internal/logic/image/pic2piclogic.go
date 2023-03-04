package image

import (
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strconv"
	"time"

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
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	// 次数限制
	today := time.Now().Format("2006-01-02")
	var total int64
	l.svcCtx.Db.Where("uid = ?", uid).Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").Count(&total)
	if total >= 3 {
		return nil, errors.New("测试阶段，每天限使用3次")
	}

	if files == nil || len(files["image"]) == 0 {
		return nil, errors.New("请上传图片")
	}

	f, err := files["image"][0].Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	hash := md5.Sum(content)
	md5String := fmt.Sprintf("%x", hash)

	pic2pic := &model.Pic2Pic{}

	l.svcCtx.Db.Where("image_hash = ?", md5String).Where("prompt = ?", prompt).Find(&pic2pic)
	if pic2pic.ID > 0 {
		return &types.ImageResponse{
			Url:  pic2pic.Url,
			Task: pic2pic.TaskId,
		}, nil
	}

	pic2pic.Uid = uint32(uid)
	pic2pic.Prompt = prompt
	pic2pic.ImageHash = md5String

	task, err := sanmuai.NewBaiduWX(l.ctx, l.svcCtx).Pic2Pic(&sanmuai.BDImageTaskRequest{
		Prompt: prompt,
		Image:  files["image"][0],
	})
	if err != nil {
		return nil, err
	}

	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/pic2pic",
		Content: strconv.Itoa(task.TaskId),
		Result:  "",
	})
	pic2pic.TaskId = strconv.Itoa(task.TaskId)
	pic2pic.AppKeyId = uint32(task.AppKeyId)
	l.svcCtx.Db.Create(&pic2pic)
	return &types.ImageResponse{
		Task: strconv.Itoa(task.TaskId),
	}, nil
}
