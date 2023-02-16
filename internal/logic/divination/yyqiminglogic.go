package divination

import (
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

type YyqimingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewYyqimingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *YyqimingLogic {
	return &YyqimingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *YyqimingLogic) Yyqiming(req *types.YYQiMingRequest, w http.ResponseWriter) (resp *types.DivinationResponse, err error) {
	other := ""
	if req.Other != "" {
		other = "，名字最好还" + req.Other
	}
	prompt := fmt.Sprintf("我中文名叫%s，性别为%s%s，请为我提供10个符合我的并且好听的英文名", req.Name, req.Sex, other)
	fmt.Println(prompt)
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           prompt,
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
				fmt.Printf("Stream error: %v\n", err)
				break
			}
			if len(response.Choices) > 0 && response.Choices[0].Text != "" {
				w.Write([]byte(response.Choices[0].Text))
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
		Type:    "divination/yyqiming",
		Content: "",
		Result:  "",
	})
	return
}
