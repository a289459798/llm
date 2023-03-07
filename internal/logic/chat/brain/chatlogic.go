package brain

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatLogic) Chat(req *types.ChatRequest, w http.ResponseWriter) (resp *types.ChatResponse, err error) {
	w.Header().Set("Content-Type", "text/event-stream")
	valid := utils.Filter(req.Message)
	if valid != "" {
		w.Write([]byte(utils.EncodeURL(valid)))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "请帮我解决一些问题",
		},
	}

	if req.ChatId != "" {
		var records []model.Record

		l.svcCtx.Db.Where("uid = ?", uid).Where("chat_id = ?", req.ChatId).Order("id asc").Find(&records)
		for _, v := range records {
			message = append(message, gogpt.ChatCompletionMessage{
				Role:    "user",
				Content: v.Content,
			})
			message = append(message, gogpt.ChatCompletionMessage{
				Role:    "assistant",
				Content: v.Result,
			})
		}

	}
	message = append(message, gogpt.ChatCompletionMessage{
		Role:    "user",
		Content: req.Message,
	})
	// 创建上下文
	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()

	ch := make(chan struct{})

	stream, err := sanmuai.NewOpenAi(ctx, l.svcCtx).CreateChatCompletionStream(message)
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	result := ""
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
				w.Write([]byte(utils.EncodeURL(response.Choices[0].Delta.Content)))
				result += response.Choices[0].Delta.Content
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
	if result == "" {
		return nil, errors.New("数据为空")
	}
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "chat/chat",
		Content: req.Message,
		Result:  result,
		ChatId:  req.ChatId,
	})

	return
}
