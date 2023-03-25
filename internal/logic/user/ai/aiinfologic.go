package ai

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AiInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiInfoLogic {
	return &AiInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiInfoLogic) AiInfo(req *types.InfoRequest) (resp *types.AIInfoResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	isVip := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db).IsVip()
	if !isVip {
		return nil, errors.New("VIP特权")
	}
	ai := &model.AI{}
	l.svcCtx.Db.Where("uid = ?", uid).Preload("Role").Find(&ai)
	if ai.ID == 0 {
		return nil, nil
	}
	return &types.AIInfoResponse{
		Name: ai.Name,
		Photo: func(img string) string {
			if img != "" {
				return fmt.Sprintf("%s?imageMogr2/thumbnail/200x200/blur/1x0/quality/75&_=%d", img, time.Now().Unix())
			}
			return ""
		}(ai.Image),
		Welcome: func(welcome string) string {
			if welcome == "" {
				return fmt.Sprintf("👋  %s你好，我是你的专属AI%s，你有任何问题都可以问我～", ai.Call, ai.Name)
			}
			return fmt.Sprintf(welcome, ai.Call, ai.Name)
		}(ai.Role.Welcome),
		Call:   ai.Call,
		Status: ai.Status,
		RoleId: ai.RoleId,
		ShowAd: true,
	}, nil
}
