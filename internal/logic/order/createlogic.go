package order

import (
	"chatgpt-tools/common/utils"
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.OrderRequest) (resp *types.OrderResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := model.AIUser{Uid: uint32(uid)}.Find(l.svcCtx.Db)
	if user.Uid == 0 {
		return nil, errors.New("用户不存在")
	}

	vip := &model.Vip{}
	l.svcCtx.Db.Where("id = ?", req.ItemId).First(&vip)
	if vip.ID == 0 {
		return nil, errors.New("vip不存在")
	}

	order := &model.Order{
		Uid:       uint32(uid),
		OrderNo:   utils.GenerateOrderNo(),
		OutNo:     fmt.Sprintf("VIP%s", utils.GenerateOrderNo()),
		OrderType: "vip",
		CostPrice: 0,
		SellPrice: vip.Price,
		PayPrice:  vip.Price,
		Status:    model.OrderStatusWaitPayment,
	}

	var orderId string
	err = l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&order).Error
		if err != nil {
			return err
		}
		err = tx.Create(&model.OrderItem{
			OrderId:   order.ID,
			ItemId:    vip.ID,
			Name:      "会员",
			Image:     "",
			CostPrice: order.CostPrice,
			SellPrice: order.SellPrice,
			PayPrice:  order.PayPrice,
			Number:    1,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.OrderResponse{
		OrderId: orderId,
		Money:   0,
	}, nil
}
