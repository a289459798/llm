package ai

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
	"io/ioutil"
	"mime/multipart"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiEditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiEditLogic {
	return &AiEditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiEditLogic) AiEdit(req *types.AIEditRequest, files map[string][]*multipart.FileHeader) (resp *types.AIInfoResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	if req.RoleId > 0 {
		template := &model.ChatTemplate{}
		l.svcCtx.Db.Where("id = ?", req.RoleId).Find(&template)
		if template.ID == 0 {
			return nil, errors.New("角色不存在")
		}
	}
	ai := &model.AI{}
	l.svcCtx.Db.Where("uid = ?", uid).Find(&ai)
	ai.Uid = uint32(uid)
	ai.Name = req.Name
	ai.Call = req.Call
	ai.RoleId = req.RoleId
	ai.Status = req.Status
	if files != nil && len(files["photo"]) > 0 {
		// 上传图片
		f, err := files["photo"][0].Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		content, err := ioutil.ReadAll(f)
		key := fmt.Sprintf("photo/uid-%d", uid)
		putPolicy := storage.PutPolicy{
			Scope: l.svcCtx.Config.Qiniu.Bucket,
		}
		mac := qbox.NewMac(l.svcCtx.Config.Qiniu.Ak, l.svcCtx.Config.Qiniu.SK)
		upToken := putPolicy.UploadToken(mac)
		cfg := storage.Config{}
		// 空间对应的机房
		cfg.Region = &storage.ZoneHuadong
		// 上传是否使用CDN上传加速
		cfg.UseCdnDomains = false
		// 构建表单上传的对象
		formUploader := storage.NewFormUploader(&cfg)
		ret := storage.PutRet{}
		dataLen := int64(len(content))
		err = formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(content), dataLen, &storage.PutExtra{})
		if err != nil {
			return nil, err
		}
		ai.Image = fmt.Sprintf("%s%s", l.svcCtx.Config.Qiniu.Domain, ret.Key)
	}

	if ai.ID == 0 {
		l.svcCtx.Db.Create(ai)
	} else {
		l.svcCtx.Db.Save(ai)
	}

	return
}
