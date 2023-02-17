package report

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
)

type WeekLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWeekLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WeekLogic {
	return &WeekLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WeekLogic) Week(req *types.ReportRequest, w http.ResponseWriter) (resp *types.ReportResponse, err error) {

	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           "请帮我把以下的工作内容填充为一篇完整的周报包含本周内容总结和下周计划,用 markdown 格式以分点叙述的形式输出:" + req.Content,
		MaxTokens:        1536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}

	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	// 创建上下文
	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()

	ch := make(chan struct{})

	stream, err := l.svcCtx.GptClient.CreateCompletionStream(ctx, gptReq)
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
		Type:    "report/week",
		Content: "",
		Result:  "",
	})
	return
}
