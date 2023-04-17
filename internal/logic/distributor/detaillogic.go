package distributor

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DistributorInfoRequest) (resp *types.DistributorInfoResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	distributor := &model.Distributor{}
	l.svcCtx.Db.Where("uid = ?", uid).Where("status = 1").Preload("Level").First(&distributor)
	if distributor.ID == 0 {
		return nil, errors.New("未开通")
	}

	return &types.DistributorInfoResponse{
		Level: distributor.Level.Name,
		Ratio: distributor.Ratio,
		Link:  fmt.Sprintf("https://chat.smuai.com/c=%d", uid),
		Money: distributor.Money,
	}, nil
}
