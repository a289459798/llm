package common

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"math/rand"
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

	consume := utils.GetSuanLi(req.Content, uint32(uid), l.svcCtx.Db)

	if amount.ChatAmount < (amount.ChatUse + uint32(consume)) {
		return nil, errors.New("次数已用完")
	}

	showAd := false

	randomNum := rand.Intn(10) + 1
	if randomNum > 6 {
		showAd = true
	}

	return &types.ValidResponse{
		Data:    strconv.Itoa(int(amount.ChatAmount) - int(amount.ChatUse)),
		ShowAd:  showAd,
		Consume: consume,
	}, nil
}
