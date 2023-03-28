package wechat

import (
	"chatgpt-tools/common/utils/appplatform"
	"chatgpt-tools/model"
	"context"
	"errors"
	"fmt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"net/http"

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

	fmt.Println(req.AppKey)
	fmt.Println("事件")
	appInfo := model.App{AppKey: req.AppKey}.Info(l.svcCtx.Db)
	if appInfo.ID == 0 {
		return nil, errors.New("App-Key 错误")
	}
	c, _ := appplatform.GetConf[appplatform.WechatOfficialConf](appInfo.Conf)
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	var cfg = &offConfig.Config{
		AppID:     c.AppId,
		AppSecret: c.AppSecret,
		Token:     c.Token,
		Cache:     memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)
	server := officialAccount.GetServer(r, w)
	fmt.Println(server.RequestMsg)
	openId := server.GetOpenID()
	fmt.Println("openId：" + openId)
	info, err := officialAccount.GetUser().GetUserInfo(openId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(info)
	return
}
