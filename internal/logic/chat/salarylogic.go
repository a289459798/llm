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
	"io"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SalaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSalaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SalaryLogic {
	return &SalaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SalaryLogic) Salary(req *types.SalaryRequest, w http.ResponseWriter) (resp *types.ChatResponse, err error) {

	w.Header().Set("Content-Type", "text/event-stream")
	valid := utils.Filter(req.Content)
	if valid != "" {
		w.Write([]byte(utils.EncodeURL(valid)))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		return
	}

	prompt := fmt.Sprintf("我需要告诉领导需要给我加薪了，我需要如何和他沟通呢？已经有一段时间没有加薪了，而且个人在岗位上也有不错的成绩，克服了一些困难，加薪后的可以更好的为公司做出更大的贡献，还有包含以下内容：%s，请用mackdown的格式输出并包含沟通的内容、技巧、以及拒绝后的方案", req.Content)

	// 创建上下文
	ctx, cancel := context.WithCancel(l.ctx)
	defer cancel()

	ch := make(chan struct{})

	stream, err := sanmuai.NewOpenAi(ctx, l.svcCtx).CreateCompletionStream(prompt)
	if err != nil {
		return nil, err
	}
	defer stream.Close()
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
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "chat/salary",
		Content: req.Content,
		Result:  "",
	})
	return
}
