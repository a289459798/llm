package divination

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QimingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQimingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QimingLogic {
	return &QimingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QimingLogic) Qiming(req *types.QiMingRequest, w http.ResponseWriter) (resp *types.DivinationResponse, err error) {
	fix := ""
	other := ""
	if req.Fix != "" {
		fix = "，名字里面必须包含" + req.Fix + "字"
	}
	if req.Other != "" {
		other = "，名字还需要符合" + req.Other
	}
	prompt := fmt.Sprintf("请帮我起一个包含%d个中文的名字，姓%s，出生年月为%s，性别为%s%s%s，请给我合适的名字以及他赋予的美好含义，请提供10个", req.Number, req.First, req.Birthday, req.Sex, fix, other)

	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	// 创建上下文
	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()

	ch := make(chan struct{})

	stream, err := sanmuai.NewOpenAi(ctx, l.svcCtx).CreateCompletionStream(prompt)
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	go func() {
		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Printf("Stream error: %v\n", err)
				break
			}
			if len(response.Choices) > 0 {
				w.Write([]byte(utils.EncodeURL(response.Choices[0].Text)))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}

		}
		close(ch)
	}()

	select {
	case <-ch:
		// 处理已完成
		logx.Infof("EventStream logic finished")
	case <-ctx.Done():
		// 处理被取消
		logx.Errorf("EventStream logic canceled")
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "divination/qiming",
		Content: "",
		Result:  "",
	})
	return
}
