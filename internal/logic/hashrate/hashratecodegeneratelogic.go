package hashrate

import (
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

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HashRateCodeGenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHashRateCodeGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HashRateCodeGenerateLogic {
	return &HashRateCodeGenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HashRateCodeGenerateLogic) HashRateCodeGenerate(req *types.HashRateCodeGenerateRequest) (resp *types.HashRateCodeGenerateResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	if uid < 1 || uid > 3 {
		return nil, errors.New("非法操作")
	}
	userId, _ := strconv.ParseInt(req.AICode, 2, 64)
	user := model.AIUser{Uid: uint32(userId)}.Find(l.svcCtx.Db)
	if user.Uid == 0 {
		return nil, errors.New("用户不存在")
	}
	str := fmt.Sprintf("%d", time.Now().Unix())
	md5sum := md5.Sum([]byte(str))
	newSign := hex.EncodeToString(md5sum[:6])
	code := strings.ToUpper(newSign)
	err = l.svcCtx.Db.Create(&model.AIHashRateCode{
		Uid:      uint32(userId),
		Code:     code,
		Day:      req.Day,
		Status:   false,
		SystemId: uint32(uid),
		AICode:   req.AICode,
		Amount:   req.Amount,
	}).Error
	if err != nil {
		return nil, err
	}
	return &types.HashRateCodeGenerateResponse{Code: code}, nil
}
