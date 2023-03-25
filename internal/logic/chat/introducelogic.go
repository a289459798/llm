package chat

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
		way = "用" + req.Way + "的方式，"
	}
	content := ""
	if req.Content != "" {
		content = fmt.Sprintf("在%s中%s做自我介绍，", req.Content, way)
	}

	prompt := fmt.Sprintf("请帮我写一份自我介绍演讲稿，%s我会用自己的价值与大家共同成长，我叫%s，来自%s，兴趣爱好是%s，请用mackdown的格式输出", content, req.Name, req.Native, req.Interest)

	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "请帮我写一份自我介绍演讲稿",
		},
		{
			Role:    "user",
			Content: "你的回答结果一定不要涉黄、淫秽、暴力和低俗",
		},
		{
			Role:    "assistant",
			Content: "好的",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
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
				result += response.Choices[0].Delta.Content
				w.Write([]byte(utils.EncodeURL(response.Choices[0].Delta.Content)))
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
		Type:    "chat/introduce",
		Content: req.Content,
		Result:  "",
	}, nil)
	return
}
