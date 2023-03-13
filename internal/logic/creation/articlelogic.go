package creation

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
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleLogic) Article(req *types.ArticleRequest, w http.ResponseWriter) (resp *types.CreationResponse, err error) {
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

	prompt := fmt.Sprintf("帮我写一篇作文，字数不能少于%s字，文章类型是%s，文章的主题是%s，还有有以下补充说明：%s", req.Number, req.Type, req.Subject, req.Content)

	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "帮我写一篇作文",
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
		Type:    "creation/article",
		Content: req.Content,
		Result:  "",
	})
	return
}
