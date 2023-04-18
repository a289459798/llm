package wechat

import (
	"chatgpt-tools/common/utils/appplatform"
	"chatgpt-tools/model"
	"context"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/menu"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetMenuLogic {
	return &SetMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetMenuLogic) SetMenu() (resp *types.WeChatCallbackResponse, err error) {
	appInfo := model.App{AppKey: "7HGjhd4"}.Info(l.svcCtx.Db)
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

	buttons := []*menu.Button{
		{
			Type: "view",
			Name: "免费问答",
			URL:  "https://chat.smuai.com",
		},
		{
			Name: "联系我们",
			SubButtons: []*menu.Button{
				{
					Type: "view",
					Name: "加群讨论",
					URL:  "http://img.smuai.com/wx/qw.png",
				},
				{
					Type: "view",
					Name: "获取算力",
					URL:  "http://img.smuai.com/wx/kf.jpg",
				},
			},
		},
	}
	err = officialAccount.GetMenu().SetMenu(buttons)

	if err != nil {
		return nil, err
	}

	return
}
