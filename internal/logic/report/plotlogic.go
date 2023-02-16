package report

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PlotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlotLogic {
	return &PlotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlotLogic) Plot(req *types.ReportRequest, w http.ResponseWriter) (resp *types.ReportResponse, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           "请帮我用以下内容完善拍摄剧本，包含故事概要和主题、人物设定、场景设定、故事结构和情节、台词和对白、拍摄风格和视觉效果、音效和音乐,用 markdown 格式以分点叙述的形式输出：" + req.Content,
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
				fmt.Println(response.Choices[0].Text)
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
	service.NewRecord(l.svcCtx.Db).Insert(model.Record{
		Uid:     l.ctx.Value("uid").(uint32),
		Type:    "report/plot",
		Content: "",
		Result:  "",
	})
	return
}
