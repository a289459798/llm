package vip

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"fmt"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VipPrivilegeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVipPrivilegeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VipPrivilegeLogic {
	return &VipPrivilegeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VipPrivilegeLogic) VipPrivilege() (resp *types.VipPrivilegeListResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := &model.AIUser{}
	l.svcCtx.Db.Preload("Vip").First(&user, uid)
	privilege := getPrivilege()

	if user.IsVip() {
		privilege.Data[0].Title = fmt.Sprintf("每天%d算力", user.Vip.Vip.Amount)
	}

	return &privilege, nil
}

func getPrivilege() types.VipPrivilegeListResponse {

	return types.VipPrivilegeListResponse{
		Data: []types.VipPrivilegeResponse{
			{Type: "suanli", Title: "每天高额算力"},
			{Type: "pic", Title: "图片高级功能"},
			{Type: "history", Title: "对话历史"},
			{Type: "ad", Title: "全站无广告"},
			{Type: "ai", Title: "专属AI"},
			{Type: "kf", Title: "专属客服"},
			{Type: "context", Title: "更多上下文"},
			{Type: "model", Title: "多模型选择"},
		},
	}
}
