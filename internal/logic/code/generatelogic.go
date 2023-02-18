package code

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

	prompt := ""
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	content := fmt.Sprintf("请用编程语言%s实现以下需求:%s，请提供代码和demo，用 maekdown 的格式输出", req.Lang, req.Content)
	if req.ChatId != "" {
		var records []model.Record

		fmt.Println(req.ChatId)
		l.svcCtx.Db.Where("uid = ?", uid).Where("chat_id = ?", req.ChatId).Order("id asc").Find(&records)
		if len(records) > 0 {
			prompt = ""
			for _, v := range records {
				prompt += "Q：" + v.Content + "\n\n"
				prompt += "A：" + v.Result + "\n\n"
			}
			content = fmt.Sprintf("%s，用%s语言", req.Content, req.Lang)
		}
	}
	prompt += "Q：" + content

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
				w.Write([]byte(utils.EncodeURL(response.Choices[0].Text)))
				result += response.Choices[0].Text
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

	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "code/generate",
		Content: content,
		Result:  result,
		ChatId:  req.ChatId,
	})
	return
}
