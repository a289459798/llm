package code

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

type GenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateLogic {
	return &GenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateLogic) Generate(req *types.GenerateRequest, w http.ResponseWriter) (resp *types.CodeResponse, err error) {
	w.Header().Set("Content-Type", "text/event-stream")
	valid := utils.Filter(req.Content, l.svcCtx.Db)
	if valid != "" {
		w.Write([]byte(utils.EncodeURL(valid)))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	content := fmt.Sprintf("请用编程语言%s实现以下需求:%s，请提供代码和demo，用 markdown 的格式输出", req.Lang, req.Content)

	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "编程问题咨询",
		},
		{
			Role:    "user",
			Content: "你的回答结果一定不要涉黄、淫秽、暴力和低俗",
		},
		{
			Role:    "assistant",
			Content: "好的",
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
		content = fmt.Sprintf("%s，用%s语言", req.Content, req.Lang)
		message = append(message, gogpt.ChatCompletionMessage{
			Role:    "user",
			Content: fmt.Sprintf("%s，用%s语言", req.Content, req.Lang),
		})
	} else {
		message = append(message, gogpt.ChatCompletionMessage{
			Role:    "user",
			Content: content,
		})
	}
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
		Type:    "code/generate",
		Content: content,
		Result:  result,
		ChatId:  req.ChatId,
	}, nil)
	return
}
