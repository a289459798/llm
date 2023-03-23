package vip

import (
	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"
	"chatgpt-tools/model"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type VipCodeGenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVipCodeGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VipCodeGenerateLogic {
	return &VipCodeGenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VipCodeGenerateLogic) VipCodeGenerate(req *types.VipCodeGenerateRequest) (resp *types.VipCodeGenerateResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	if uid < 1 || uid > 3 {
		return nil, errors.New("非法操作")
	}
	userId, _ := strconv.ParseInt(req.AICode, 2, 64)
	user := model.User{ID: uint32(userId)}.Find(l.svcCtx.Db)
	if user.ID == 0 {
		return nil, errors.New("用户不存在")
	}
	str := fmt.Sprintf("%d", time.Now().Unix())
	md5sum := md5.Sum([]byte(str))
	newSign := hex.EncodeToString(md5sum[:6])
	code := strings.ToUpper(newSign)
	err = l.svcCtx.Db.Create(&model.VipCode{
		Uid:    uint32(uid),
		Code:   code,
		VipId:  req.VipId,
		Day:    req.Day,
		Status: false,
	}).Error
	if err != nil {
		return nil, err
	}
	return &types.VipCodeGenerateResponse{Code: code}, nil
}
