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
	w.Header().Set("Content-Type", "text/event-stream")
	valid := utils.Filter(req.Content, l.svcCtx.Db)
	if valid != "" {
		w.Write([]byte(utils.EncodeURL(valid)))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()

	ch := make(chan struct{})

	prompt := fmt.Sprintf("请帮我完善一份策划方案，活动类型是%s，活动主要目的%s，时间周期为%s，主要针对%s，以下是主要活动内容：%s，需要提供完整的活动方案，包括但不限于前期准备、活动的实施方案、活动过程跟踪、效果不及预期的方案、活动效果、需要的支持等，请用mackdown的格式输出", req.Way, req.Target, req.Date, req.User, req.Content)

	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "帮我写一份活动策划",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

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
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "creation/activity",
		Content: req.Content,
		Result:  "",
	}, nil)
	return
}
