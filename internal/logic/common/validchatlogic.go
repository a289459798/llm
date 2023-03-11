package common

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidChatLogic {
	return &ValidChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidChatLogic) ValidChat(req *types.ValidRequest) (resp *types.ValidResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	amount := model.NewAccount(l.svcCtx.Db).GetAccount(uint32(uid), time.Now())

	consume := 1
	if req.Content == "image/create" {
		consume = 3
	}

	if amount.ChatAmount < (amount.ChatUse + uint32(consume)) {
		return nil, errors.New("次数已用完")
	}

	showAd := false
	var total int64
	l.svcCtx.Db.Model(&model.Record{}).Where("uid = ?", uid).Where("type = ?", req.Content).Count(&total)
	if total > 0 && (total%5 == 0) {
		showAd = true
	}
	return &types.ValidResponse{
		Data:    strconv.Itoa(int(amount.ChatAmount) - int(amount.ChatUse)),
		ShowAd:  showAd,
		Consume: consume,
	}, nil
}
