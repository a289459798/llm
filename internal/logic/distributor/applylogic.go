package distributor

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyLogic {
	return &ApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyLogic) Apply(req *types.DistributorApplyRequest) (resp *types.DistributorApplyResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	distributor := &model.Distributor{}
	l.svcCtx.Db.Where("uid = ?", uid).First(&distributor)
	if distributor.ID > 0 {
		return nil, errors.New("请勿重复开通")
	}
	distributor.Uid = uint32(uid)
	level := &model.DistributorLevel{}
	l.svcCtx.Db.Order("id asc").First(level)
	distributor.LevelId = level.ID
	distributor.Ratio = level.Ratio
	distributor.Status = true
	l.svcCtx.Db.Create(distributor)
	return
}
