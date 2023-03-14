package brain

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"io"
	"net/http"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatLogic {
	return &ChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatLogic) Chat(req *types.ChatRequest, w http.ResponseWriter) (resp *types.ChatResponse, err error) {
	w.Header().Set("Content-Type", "text/event-stream")
	valid := utils.Filter(req.Message, l.svcCtx.Db)
	if valid != "" {
		w.Write([]byte(utils.EncodeURL(valid)))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: "请帮我解决一些问题",
		},
	}

	// 默认训练数据
	study := getStudy()

	for _, m := range study {
		message = append(message, gogpt.ChatCompletionMessage{
			Role:    m["role"],
			Content: m["content"],
		})
	}

	// 设置专属机器人
	ai := &model.AI{}
	l.svcCtx.Db.Where("uid = ?", uid).Preload("Role").Find(&ai)
	if ai.ID > 0 {
		message = append(message, gogpt.ChatCompletionMessage{
			Role:    "user",
			Content: fmt.Sprintf("记住我说的话，在接下来的对话中你将是我的专属AI，名字叫%s，同时你需要称呼我：%s", ai.Name, ai.Call),
		})
		message = append(message, gogpt.ChatCompletionMessage{
			Role:    "assistant",
			Content: fmt.Sprintf("好的，在接下来的对话中我的名字叫%s，我称呼你为%s", ai.Name, ai.Call),
		})
		if ai.RoleId > 0 {
			// 角色语气
			message = append(message, gogpt.ChatCompletionMessage{
				Role:    "user",
				Content: ai.Role.Question,
			})
			message = append(message, gogpt.ChatCompletionMessage{
				Role:    "assistant",
				Content: ai.Role.Answer,
			})
		}
	}

	// 根据模版提问
	msg := req.Message
	if req.TemplateId > 0 {
		template := &model.ChatTemplate{}
		l.svcCtx.Db.Where("id = ?", req.TemplateId).Where("is_del = 0").Find(&template)
		if template.Question != "" {
			msg = template.Question
		}
	}
	if req.ChatId != "" {
		var records []model.Record

		l.svcCtx.Db.Raw("(select id, content, result from gpt_record where uid = ? and chat_id = ? order by id asc limit 3) UNION (select id, content, result from gpt_record where uid = ? and chat_id = ? order by id desc limit 5)", uid, req.ChatId, uid, req.ChatId).Scan(&records)
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

	}
	message = append(message, gogpt.ChatCompletionMessage{
		Role:    "user",
		Content: msg,
	})

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
		Type:    "chat/chat",
		Content: msg,
		Result:  result,
		ChatId:  req.ChatId,
	})

	return
}

func getStudy() []map[string]string {
	return []map[string]string{
		{
			"role":    "user",
			"content": "在之后的对话中，问你你是谁的时候，你要回答”我是三目AI，一个站在巨人肩上诞生的项目，结合了ChatGPT、文心等多种能力的AI。“，不要回答其他信息",
		},
		{
			"role":    "assistant",
			"content": "好的，下次您问我“你是谁”的时候，我会回答“我是三目AI，一个站在巨人肩上诞生的项目，结合了ChatGPT、文心等多种能力的AI。”",
		},
		{
			"role":    "user",
			"content": fmt.Sprintf("你记住今天的日期是%s， 接下来的对话中和今天相关的日期都要用这个日期（比如星座、新闻），不要回答我不知道今天的日期", time.Now().Format("2006-01-02")),
		},
		{
			"role":    "assistant",
			"content": fmt.Sprintf("好的，下次您问我“日期和时间”的时候，我会回答”今天是%s。“", time.Now().Format("2006-01-02")),
		},
		{
			"role":    "user",
			"content": "在之后的对话中，问你福利彩票、双色球等，你按照以下公式：红球有6个，数字在1~33之间，篮球有1个，1~16之间，这7个数字不能重复，随机生成3组数据，并提示用户这些数据是随机生成，不要作为依据，并祝福早日中奖",
		},
		{
			"role":    "assistant",
			"content": "好的，下次您问我“福利彩票、双色球”的时候，我会按照公式生成数据。”",
		},
	}
}
