package task

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskShareFollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskShareFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskShareFollowLogic {
	return &TaskShareFollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskShareFollowLogic) TaskShareFollow(req *types.TaskShareFollowRequest) (resp *types.TaskResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	user := &model.AIUser{}
	l.svcCtx.Db.Where("open_id = ?", req.OpenId).Where("uid != ?", uid).Find(&user)
	if user.Uid == 0 {
		return nil, errors.New("用户不存在")
	}

	today := time.Now().Format("2006-01-02")
	var total int64
	l.svcCtx.Db.Model(&model.ShareRecord{}).
		Where("uid = ?", user.Uid).
		Where("follow_id = ?", uid).
		Count(&total)
	if total > 0 {
		return nil, errors.New("重复推荐")
	}

	taskType := "share_follow"
	l.svcCtx.Db.Model(&model.AccountRecord{}).
		Where("uid = ?", user.Uid).
		Where("type = ?", taskType).
		Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").
		Count(&total)

	if total > 3 {
		return nil, errors.New("已超过最大奖励次数")
	}

	var add uint32 = 5

	l.svcCtx.Db.Transaction(func(tx *gorm.DB) error {
		// 插入
		tx.Create(&model.ShareRecord{
			Uid:      user.Uid,
			FollowId: uint32(uid),
		})
		// 增加次数
		amount := model.NewAccount(tx).GetAccount(user.Uid, time.Now())
		amount.ChatAmount += add
		tx.Save(&amount)

		// 记录
		tx.Create(&model.AccountRecord{
			Uid:           user.Uid,
			RecordId:      0,
			Way:           1,
			Type:          taskType,
			Amount:        add,
			CurrentAmount: amount.ChatAmount - amount.ChatUse,
		})
		return nil
	})

	return
}
