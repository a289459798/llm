package chat

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

type IntroduceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIntroduceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IntroduceLogic {
	return &IntroduceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IntroduceLogic) Introduce(req *types.IntroduceRequest, w http.ResponseWriter) (resp *types.ChatResponse, err error) {

	way := ""
	if req.Way != "" {
		way = "介绍方式需要" + req.Way
	}
	content := ""
	if req.Content != "" {
		content = "还需要包含以下信息：" + req.Content
	}

	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           fmt.Sprintf("请帮我写一份自我介绍演讲稿，我叫%s，来自%s，兴趣爱好是%s%s%s，请用mackdown的格式输出", req.Name, req.Native, req.Interest, way, content),
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
		Type:    "chat/introduce",
		Content: "",
		Result:  "",
	})
	return
}