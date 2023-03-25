package common

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/security"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidTextLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidTextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidTextLogic {
	return &ValidTextLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidTextLogic) ValidText(req *types.ValidRequest) (resp *types.ValidResponse, err error) {
	valid := utils.Filter(req.Content, l.svcCtx.Db)
	if valid != "" {
		return nil, errors.New(valid)
	}
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := &model.AIUser{}
	l.svcCtx.Db.Where("uid = ?", uid).First(user)
	if user.Uid == 0 {
		return nil, errors.New("用户不存在")
	}

	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &config.Config{
		AppID:     l.svcCtx.Config.MiniApp.AppId,
		AppSecret: l.svcCtx.Config.MiniApp.AppSecret,
		Cache:     memory,
	}
	mini := wc.GetMiniProgram(cfg)
	res, _ := mini.GetSecurity().MsgCheck(&security.MsgCheckRequest{
		OpenID:  user.OpenId,
		Scene:   3,
		Content: req.Content,
	})

	if res.Result.Suggest != "pass" {
		resMessage := ""
		for _, s := range res.Detail {
			if s.Suggest != "pass" {
				if s.Prob > 80 {
					resMessage += fmt.Sprintf("%s:%s ", s.Label, s.Keyword)
				}
			}
		}
		if resMessage != "" {
			return nil, errors.New(fmt.Sprintf("包含：%s", resMessage))
		}
	}

	return
}
