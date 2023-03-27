package image

import (
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/model"
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskLogic {
	return &TaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskLogic) Task(req *types.ImageTaskRequest) (resp *types.ImageMultiAsyncResponse, err error) {
	ai := sanmuai.GetAI("Tencentarc", sanmuai.SanmuData{
		Ctx:    l.ctx,
		SvcCtx: l.svcCtx,
	})

	stream, err := ai.ImageTask(sanmuai.ImageAsyncTask{
		Task: req.Task,
	})
	if err != nil {
		return nil, err
	}
	if len(stream.Output) > 0 {
		uid, _ := l.ctx.Value("uid").(json.Number).Int64()
		record := &model.Record{}
		l.svcCtx.Db.Where("uid = ?", uid).
			Where("type = ?", "image/pic-repair").
			Where("chat_id = ?", req.Task).
			First(&record)
		if record.ID > 0 {
			record.Result = strings.Join(stream.Output, ",")
			l.svcCtx.Db.Save(&record)
		}
		for i := 0; i < len(stream.Output); i++ {
			stream.Output[i] = base64.StdEncoding.EncodeToString([]byte(stream.Output[i]))
		}
	}

	return &types.ImageMultiAsyncResponse{
		Url: stream.Output,
	}, nil
}
