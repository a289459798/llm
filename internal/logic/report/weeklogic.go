package report

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
	gogpt "github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
)

type WeekLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWeekLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WeekLogic {
	return &WeekLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WeekLogic) Week(req *types.ReportRequest, w http.ResponseWriter) (resp *types.ReportResponse, err error) {

	w.Header().Set("Content-Type", "text/event-stream")
	valid := utils.Filter(req.Content, l.svcCtx.Db)
	if valid != "" {
		w.Write([]byte(utils.EncodeURL(valid)))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return
	}

	prompt := "请帮我把以下的工作内容填充为一篇完整的周报包含本周内容、下周计划、本周总结,用 markdown 格式以分点叙述的形式输出:" + req.Content

	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "帮我完善一下周报",
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
		Type:    "report/week",
		Content: req.Content,
		Result:  "",
	})
	return
}
