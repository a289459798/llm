package game

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/common/utils/sanmuai"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type IdiomAnswerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIdiomAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IdiomAnswerLogic {
	return &IdiomAnswerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IdiomAnswerLogic) IdiomAnswer(req *types.IdiomRequest) (resp *types.GameResponse, err error) {
	rc := []rune(req.Content)
	if req.Pre != "" {
		rp := []rune(req.Pre)
		if len(rc) != 4 {
			return nil, errors.New("请回答四字成语")
		}
		if string(rc[0]) != string(rp[3]) {
			return nil, errors.New("请回答正确的四字成语")
		}
		//gptReq := gogpt.CompletionRequest{
		//	Model:            gogpt.GPT3TextDavinci003,
		//	Prompt:           fmt.Sprintf("请告诉我%s是成语吗？只要回答我对还是错，不要包含其他文字与字符", req.Content),
		//	MaxTokens:        1536,
		//	Temperature:      0.7,
		//	TopP:             1,
		//	FrequencyPenalty: 0,
		//	PresencePenalty:  0,
		//	N:                1,
		//}
		//stream, err := l.svcCtx.GptClient.CreateCompletion(l.ctx, gptReq)
		//if err != nil {
		//	return nil, err
		//}
		//
		//if utils.TrimHtml(stream.Choices[0].Text) == "错" {
		//	return nil, errors.New("这好像不是一个成语")
		//}
	}

	prompt := fmt.Sprintf("请给我一个以%s开头的正确的四个字成语，不要包含句号在内的所有字符", string(rc[3]))

	stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateCompletion(prompt)
	if err != nil {
		return nil, err
	}
	return &types.GameResponse{Data: utils.TrimHtml(stream.Choices[0].Text)}, nil
}
