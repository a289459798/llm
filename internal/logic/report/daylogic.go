package report

import (
	"chatgpt-tools/common/utils"
	"context"
	"errors"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DayLogic {
	return &DayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DayLogic) Day(req *types.ReportRequest, w http.ResponseWriter) (resp *types.ReportResponse, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           "请帮我把以下的工作内容填充为一篇完整的日报,用 markdown 格式以分点叙述的形式输出:" + req.Content,
		MaxTokens:        1536,
		Temperature:      0.7,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		N:                1,
	}
	w.Header().Set("Content-Type", "text/event-stream")
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
	return
}
