package efficiency

import (
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

	"github.com/zeromicro/go-zero/core/logx"
)

type MindLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MindLogic {
	return &MindLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MindLogic) Mind(req *types.MindRequest) (resp *types.EfficiencyResponse, err error) {
	tools := model.ToolsMind
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db)
	message, isFirst, err := model.Record{Uid: uint32(uid), ChatId: req.ChatId, Type: tools}.GetMessage(l.svcCtx.Db, user)
	if err != nil {
		return nil, err
	}

	title := ""
	if isFirst {
		title = req.Content
	}
	showContent := req.Content
	content := fmt.Sprintf("结合上下文把以下内容:%s，转换成思维导图的json格式", req.Content)

	message = append(message, gogpt.ChatCompletionMessage{
		Role:    gogpt.ChatMessageRoleUser,
		Content: content,
	})

	// 创建上下文

	stream, err := sanmuai.NewOpenAi(l.ctx, l.svcCtx).CreateChatCompletion(message)
	fmt.Println(stream)
	if err != nil {
		return nil, err
	}
	result := stream.Choices[0].Message.Content

	if result == "" {
		return nil, errors.New("数据为空")
	}
	service.NewRecord(l.svcCtx.Db).Insert(&model.Record{
		Uid:         uint32(uid),
		Type:        tools,
		Title:       title,
		Content:     content,
		ShowContent: showContent,
		ChatId:      req.ChatId,
		Result:      result,
	}, nil)
	return &types.EfficiencyResponse{
		Data: result,
	}, nil
}
