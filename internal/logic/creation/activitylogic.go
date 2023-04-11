package creation

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

type ActivityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActivityLogic {
	return &ActivityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActivityLogic) Activity(req *types.ActivityRequest, w http.ResponseWriter) (resp *types.CreationResponse, err error) {

	tools := model.ToolsActivity
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
		showContent = fmt.Sprintf("活动类型：%s\n活动目的：%s\n活动时间：%s\n用户：%s\n说明：%s", req.Way, req.Target, req.Date, req.User, req.Content)
		title = showContent
		content = fmt.Sprintf("活动类型是%s，活动主要目的%s，时间周期为%s，主要针对%s，以下是主要活动内容：%s", req.Way, req.Target, req.Date, req.User, req.Content)
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
