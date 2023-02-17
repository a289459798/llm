package divination

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type JiemengLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJiemengLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JiemengLogic {
	return &JiemengLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JiemengLogic) Jiemeng(req *types.JieMengRequest, w http.ResponseWriter) (resp *types.DivinationResponse, err error) {
	gptReq := gogpt.CompletionRequest{
		Model:            gogpt.GPT3TextDavinci003,
		Prompt:           "我昨天晚上做了一个梦，请帮我解一下这个梦的详细含义提供更多的信息，请用markdown格式输出，以下是做梦的大致内容:" + req.Content,
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
			if len(response.Choices) > 0 {
				w.Write([]byte(utils.EncodeURL(response.Choices[0].Text)))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}

		}
		w.Write([]byte(utils.EncodeURL("\n\n***\n\n*请注意，解梦是一种主观的解释和理解，其准确性和可信度存在较大的差异。解梦结果仅供参考，不应作为决策或行动的依据。*")))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
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
		Type:    "divination/jiemeng",
		Content: "",
		Result:  "",
	})
	return
}
