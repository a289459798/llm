package game

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"context"
	gogpt "github.com/sashabaranov/go-openai"
	"github.com/zeromicro/go-zero/core/logx"
)

type IdiomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdiomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdiomLogic {
	return &IdiomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdiomLogic) Idiom() (resp *types.GameResponse, err error) {

	message := []gogpt.ChatCompletionMessage{}
	message = append(message, gogpt.ChatCompletionMessage{
		Role:    "system",
		Content: "生成成语",
	})
	message = append(message, gogpt.ChatCompletionMessage{
		Role:    "user",
		Content: "请随机给我生成一个四字成语，不要包含句号在内的所有字符",
	})
	stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateChatCompletion(message)
	if err != nil {
		return nil, err
	}
	return &types.GameResponse{Data: utils.TrimHtml(stream.Choices[0].Message.Content)}, nil
}
