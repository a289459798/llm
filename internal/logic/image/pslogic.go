package image

import (
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/model"
	"chatgpt-tools/service"
	"context"
	"encoding/json"
	"fmt"
	gogpt "github.com/sashabaranov/go-openai"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PsLogic {
	return &PsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PsLogic) Ps(req *types.PSRequest) (resp *types.ImageMultiAsyncResponse, err error) {
	ai := sanmuai.GetAI("Paintbytext", sanmuai.SanmuData{
		Ctx:    l.ctx,
		SvcCtx: l.svcCtx,
	})

	content := req.Content
	message := []gogpt.ChatCompletionMessage{
		{
			Role:    "user",
			Content: fmt.Sprintf("现在起你化身为翻译官，把我提供内容翻译成英文，引号里面的不要翻译,也不需要任何解释，比如我说狗，你回复\"dog\"，第一句要翻译的内容是：%s", req.Content),
		},
	}
	stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateChatCompletion(message)
	if err != nil {
		return nil, err
	}
	if len(stream.Choices) > 0 && stream.Choices[0].Message.Content != "" {
		content = stream.Choices[0].Message.Content
	}

	result, err := ai.ImagePSAsync(sanmuai.ImagePS{Image: req.Image, Text: content})
	if err != nil {
		return nil, err
	}

	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:     uint32(uid),
		Type:    "image/ps",
		Content: req.Image,
		ChatId:  result.Task,
		Model:   "Paintbytext",
	}, nil)

	return &types.ImageMultiAsyncResponse{
		Model: "Tencentarc",
		Task:  result.Task,
	}, nil
}
