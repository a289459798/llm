package hashrate

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HashRateCxchangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHashRateCxchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HashRateCxchangeLogic {
	return &HashRateCxchangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HashRateCxchangeLogic) HashRateCxchange(req *types.HashRateCxchangeRequest) (resp *types.HashRateCxchangeResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	hashRateCode := &model.AIHashRateCode{}
	l.svcCtx.Db.Where("uid = ?", uid).Where("code = ?", req.Code).Where("status != 1").First(&hashRateCode)
	if hashRateCode.ID == 0 {
		return nil, errors.New("兑换码错误")
	}

	tx := l.svcCtx.Db.Begin()
	hashRateCode.Status = true
	tx.Save(&hashRateCode)
	err = tx.Create(&model.AIUserHashRate{
		Uid:       hashRateCode.Uid,
		Amount:    hashRateCode.Amount,
		UseAmount: 0,
		Expiry:    time.Now().AddDate(0, 0, int(hashRateCode.Day)),
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("兑换错误")
	}
	account := model.NewAccount(tx).GetAccount(uint32(uid), time.Now())
	err = tx.Create(&model.AccountRecord{
		Uid:           hashRateCode.Uid,
		RecordId:      0,
		Way:           1,
		Type:          "exchange",
		Amount:        hashRateCode.Amount,
		CurrentAmount: account.Amount,
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("兑换错误")
	}
	tx.Commit()

	return &types.HashRateCxchangeResponse{}, nil
}
