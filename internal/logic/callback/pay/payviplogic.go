package pay

import (
	"chatgpt-tools/common/utils/pay"
	"context"
	"fmt"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayVipLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayVipLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayVipLogic {
	return &PayVipLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayVipLogic) PayVip(req *types.PayRequest, r *http.Request) (resp *types.WechatPayResponse, err error) {
	payModel := pay.GetPay(req.Type, pay.PayData{
		Ctx:      l.ctx,
		Config:   l.svcCtx.Config,
		Merchant: req.Merchant,
	})

	payNotify, err := payModel.PayNotify(r)
	if err != nil {
		return nil, err
	}
	fmt.Println(payNotify)

	return &types.WechatPayResponse{
		Data: "success",
	}, nil
}
