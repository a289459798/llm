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
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DiaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDiaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DiaryLogic {
	return &DiaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DiaryLogic) Diary(req *types.DiaryRequest, w http.ResponseWriter) (resp *types.CreationResponse, err error) {
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

	week := int(time.Now().Weekday())
	weekDay := "未知"
	switch week {
	case 0:
		weekDay = "日"
		break
	case 1:
		weekDay = "一"
		break
	case 2:
		weekDay = "二"
		break
	case 3:
		weekDay = "三"
		break
	case 4:
		weekDay = "四"
		break
	case 5:
		weekDay = "五"
		break
	case 6:
		weekDay = "六"
		break
	}
	prompt := fmt.Sprintf("今天是%s，星期%s，%s的天气%s，大致内容是：%s，请帮我写一份完整了日记，字数不能少于300字", time.Now().Format("2006-01-02"), weekDay, req.City, req.Weather, req.Content)

	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "帮我写一篇日记",
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
	})
	return
}
