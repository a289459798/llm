package wechat

import (
	"context"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscribeCallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscribeCallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeCallLogic {
	return &SubscribeCallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscribeCallLogic) SubscribeCall(r *http.Request, w http.ResponseWriter) (resp *types.WeChatCallbackResponse, err error) {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     l.svcCtx.Config.OfficialAccount.AppId,
		AppSecret: l.svcCtx.Config.OfficialAccount.AppSecret,
		Cache:     memory,
		Token:     l.svcCtx.Config.OfficialAccount.Token,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	// 传入request和responseWriter
	server := officialAccount.GetServer(r, w)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err = server.Serve()
	if err != nil {
		return nil, err
	}
	//发送回复的消息
	server.Send()

	return
}
