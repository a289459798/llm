package image

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type Image2TextLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImage2TextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Image2TextLogic {
	return &Image2TextLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Image2TextLogic) Image2Text(req *types.Image2TextRequest, w http.ResponseWriter) (resp *types.Image2TextResponse, err error) {
	w.Header().Set("Content-Type", "text/event-stream")
	ai := sanmuai.GetAI("Salesforce", sanmuai.SanmuData{
		Ctx:    l.ctx,
		SvcCtx: l.svcCtx,
	})

	imageText, err := ai.ImageText(sanmuai.Image2Text{Image: req.Image})
	if err != nil {
		return nil, err
	}
	prompt := fmt.Sprintf("一张图片包含以下内容：%s，请帮我组织一下语言形成一篇短文，用中文输出", imageText)

	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "看图写话",
		},
		{
			Role:    "user",
			Content: prompt,
		},
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
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/img2text",
		Content: req.Image,
		Result:  result,
		Model:   "Salesforce",
	}, nil)
	return

}
