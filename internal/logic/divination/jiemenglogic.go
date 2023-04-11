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
	gogpt "github.com/sashabaranov/go-openai"
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
	tools := model.ToolsJieMeng
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db)
	message, isFirst, err := model.Record{Uid: uint32(uid), ChatId: req.ChatId, Type: tools}.GetMessage(l.svcCtx.Db, user)
	if err != nil {
		return nil, err
	}

	content := req.Content
	showContent := ""
	title := ""
	if isFirst {
		title = req.Content
		showContent = req.Content
		content = fmt.Sprintf("我梦里面的内容是：%s", content)
	}

	message = append(message, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleUser,
		Content: content,
	})

	// 创建上下文
	w.Header().Set("Content-Type", "text/event-stream")
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
				fmt.Printf("Stream error: %v\n", err)
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
	if result == "" {
		return nil, errors.New("数据为空")
	}
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:         uint32(uid),
		Type:        tools,
		Title:       title,
		Content:     content,
		ShowContent: showContent,
		ChatId:      req.ChatId,
		Result:      result,
	}, nil)
	return
}
