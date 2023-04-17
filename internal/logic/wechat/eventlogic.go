package wechat

import (
	"chatgpt-tools/common/utils/appplatform"
	"chatgpt-tools/model"
	"context"
	"errors"
	"fmt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"net/http"
	"strings"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EventLogic {
	return &EventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventLogic) Event(req types.WechatValidateRequest, r *http.Request, w http.ResponseWriter) (resp *types.WeChatCallbackResponse, err error) {

	appInfo := model.App{AppKey: req.AppKey}.Info(l.svcCtx.Db)
	if appInfo.ID == 0 {
		return nil, errors.New("App-Key 错误")
	}
	c, _ := appplatform.GetConf[appplatform.WechatOfficialConf](appInfo.Conf)
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	var cfg = &offConfig.Config{
		AppID:          c.AppId,
		AppSecret:      c.AppSecret,
		Token:          c.Token,
		EncodingAESKey: c.EncodingAESKey,
		Cache:          memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	server := officialAccount.GetServer(r, w)
	// 设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {

		switch server.RequestMsg.MsgType {
		case message.MsgTypeEvent:
			switch server.RequestMsg.Event {
			case message.EventSubscribe:
				text := message.NewText("欢迎关注三目")
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
			}
			break
		}
		return nil
	})

	err = server.Serve()
	if err != nil {
		return nil, err
	}
	openId := server.GetOpenID()
	switch server.RequestMsg.MsgType {
	case message.MsgTypeEvent:
		switch server.RequestMsg.Event {
		case message.EventScan:
			if strings.Index(server.RequestMsg.EventKey, "login_") >= 0 {
				err = l.login(server.RequestMsg.EventKey, openId, req.AppKey, officialAccount)
				if err != nil {
					return nil, err
				}
			}
			break

		case message.EventSubscribe:
			if strings.Index(server.RequestMsg.EventKey, "qrscene_login") >= 0 {
				server.RequestMsg.EventKey = strings.Replace(server.RequestMsg.EventKey, "qrscene_", "", 1)
				err = l.login(server.RequestMsg.EventKey, openId, req.AppKey, officialAccount)
				if err != nil {
					return nil, err
				}
			}
			break

		}
		break
	}
	server.Send()
	return
}

func (l *EventLogic) login(key string, openId string, appKey string, officialAccount *officialaccount.OfficialAccount) error {
	scan := &model.ScanScene{}
	l.svcCtx.Db.Where("scene = ?", key).First(&scan)
	if scan.ID > 0 {
		info, err := officialAccount.GetUser().GetUserInfo(openId)
		if err != nil {
			return err
		}
		aiUser, tokenString, err := model.AIUser{}.Login(l.svcCtx.Db, model.UserLogin{
			OpenID:       openId,
			UnionID:      info.UnionID,
			Channel:      scan.Channel,
			AppKey:       appKey,
			AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
			AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		})

		if err != nil {
			return err
		}
		scan.Data = fmt.Sprintf("%d|%s", aiUser.Uid, tokenString)
		l.svcCtx.Db.Save(&scan)

		_, err = officialAccount.GetTemplate().Send(&message.TemplateMessage{
			ToUser:     info.OpenID,
			TemplateID: "sxKNKCHFXS9-r-WAM0Ay1DFTF4SUH0EeC8SrNOao3sw",
			Data: map[string]*message.TemplateDataItem{
				"first": {
					Value: "登录成功",
				},
				"keyword1": {
					Value: fmt.Sprintf("%b", aiUser.Uid),
				},
				"keyword2": {
					Value: "三目",
				},
				"keyword3": {
					Value: time.Now().Format("2006-01-02 15:04:05"),
				},
			},
		})
		if err != nil {
			fmt.Println("TemplateMessage err")
			fmt.Println(err)
		}
	}
	return nil
}
