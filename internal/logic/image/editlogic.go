package image

import (
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"io/ioutil"
	"mime/multipart"
)

type EditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditLogic {
	return &EditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditLogic) Edit(req *types.ImageEditRequest, files map[string][]*multipart.FileHeader) (resp *types.ImageResponse, err error) {
	prompt := req.Content
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()

	if files == nil || len(files["image"]) == 0 {
		return nil, errors.New("请上传图片")
	}

	f, err := files["image"][0].Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file, err := ioutil.TempFile("", "upload")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// 将文件内容复制到临时文件中
	_, err = io.Copy(file, f)
	if err != nil {
		return nil, err
	}
	stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateEditImage(file, prompt)
	if err != nil {
		return nil, err
	}

	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/edit",
		Content: req.Content,
		Result:  "",
	}, nil)

	return &types.ImageResponse{
		Url: stream.Data[0].B64JSON,
	}, nil
}
