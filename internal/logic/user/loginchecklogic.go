package user

import (
	"chatgpt-tools/model"
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginCheckLogic {
	return &LoginCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginCheckLogic) LoginCheck(req *types.LoginCheckRequest) (resp *types.InfoResponse, err error) {
	data := &model.ScanScene{}
	l.svcCtx.Db.Where("scene = ?", req.SceneStr).Where("data != ?", "").First(&data)
	if data.ID == 0 {
		return nil, errors.New(string(http.StatusNotFound))
	}

	str := strings.Split(data.Data, "|")
	uid, _ := strconv.Atoi(str[0])
	return &types.InfoResponse{
		Token: str[1],
		Uid:   uint32(uid),
	}, nil
}
