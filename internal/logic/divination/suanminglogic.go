package divination

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SuanmingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSuanmingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SuanmingLogic {
	return &SuanmingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SuanmingLogic) Suanming(req *types.SuanMingRequest, w http.ResponseWriter) (resp *types.DivinationResponse, err error) {
	sex := ""
	content := ""
	if req.Sex != "" {
		sex = "，性别为" + req.Sex
	}
	if req.Content != "" {
		content = "，还有以下内容参考：" + req.Content
	}
	prompt := fmt.Sprintf("请帮我算命，我叫%s，出生年月为%s%s%s，给一份详情的算命报告，包含八字分析、五行分析、命理分析、事业分析、爱情分析、财运分析等相关内容，请用markdown格式输出", req.Name, req.Birthday, sex, content)

	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	// 创建上下文
	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()

	ch := make(chan struct{})

	stream, err := sanmuai.NewOpenAi(ctx, l.svcCtx).CreateCompletionStream(prompt)
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
				fmt.Printf("Stream error: %v\n", err)
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
		w.Write([]byte(utils.EncodeURL("\n\n***\n\n*请注意，算命是一种非科学的预测方法，其准确性和可信度存在较大的差异。算命结果仅供参考，不应作为决策或行动的依据。*")))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
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
		Type:    "divination/suanming",
		Content: "",
		Result:  "",
	}, nil)
	return
}
