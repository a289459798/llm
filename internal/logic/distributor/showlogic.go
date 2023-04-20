package distributor

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.DistributorShowRequest) (resp *types.DistributorShowResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	distributor := &model.Distributor{}
	l.svcCtx.Db.Where("uid = ?", uid).First(&distributor)
	if distributor.ID > 0 {
		return &types.DistributorShowResponse{Show: false}, nil
	}
	return &types.DistributorShowResponse{Show: true}, nil
}
