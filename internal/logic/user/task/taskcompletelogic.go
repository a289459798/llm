package task

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskCompleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskCompleteLogic {
	return &TaskCompleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskCompleteLogic) TaskComplete(req *types.TaskRequest, r *http.Request) (resp *types.TaskCompleteResponse, err error) {
	timestamp := r.Header.Get("timestamp")
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	accountRecord := &model.AccountRecord{}
	l.svcCtx.Db.Where("uid = ?", uid).
		Where("record_id = ?", timestamp).
		Find(&accountRecord)
	if accountRecord.ID > 0 {
		return nil, errors.New("非法请求")
	}

	user := &model.AIUser{}
	l.svcCtx.Db.Where("uid = ?", uid).First(&user)
	if user.Uid == 0 {
		return nil, errors.New("用户不存在")
	}

	if req.Type == "group" && user.JoinGroup {
		return nil, errors.New("重复操作")
	}

	today := time.Now().Format("2006-01-02")
	var total int64
	l.svcCtx.Db.Model(&model.AccountRecord{}).
		Where("uid = ?", uid).
		Where("type = ?", req.Type).
		Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").
		Count(&total)

	var add uint32 = 5
	if (req.Type == "share" && total >= 3) || (req.Type == "ad" && total >= 10) {
		return nil, errors.New("已超过最大任务次数")
	} else if req.Type == "ad" && total == 0 {
		add = 10
	} else if req.Type == "group" {
		add = 10
	}

	tx := l.svcCtx.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 增加次数
	amount := model.NewAccount(tx).GetAccount(uint32(uid), time.Now())
	amount.ChatAmount += add
	tx.Save(&amount)

	if req.Type == "group" {
		user.JoinGroup = true
		l.svcCtx.Db.Save(user)
	}

	// 记录
	t, _ := strconv.Atoi(timestamp)
	err = tx.Create(&model.AccountRecord{
		Uid:           uint32(uid),
		RecordId:      uint32(t),
		Way:           1,
		Type:          req.Type,
		Amount:        add,
		CurrentAmount: amount.ChatAmount - amount.ChatUse,
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("重复执行")
	}
	var welfare uint32 = 0

	tasksSetting, err := model.Setting{Name: "task_welfare"}.Find(l.svcCtx.Db)
	if err == nil && len(tasksSetting) > 0 {
		startTime, err1 := time.ParseInLocation("2006-01-02 15:04:05", tasksSetting["start"].(string)+" 00:00:00", time.Local)
		endTime, err2 := time.ParseInLocation("2006-01-02 15:04:05", tasksSetting["end"].(string)+" 23:59:59", time.Local)
		if err1 == nil && err2 == nil && time.Now().Unix() >= startTime.Unix() && time.Now().Unix() < endTime.Unix() {
			l.svcCtx.Db.Model(&model.AccountRecord{}).
				Where("uid = ?", uid).
				Where("type in (?, ?)", "ad", "share").
				Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").
				Count(&total)

			if total == 12 {
				welfare = uint32(tasksSetting["amount"].(float64))
				amount := model.NewAccount(tx).GetAccount(uint32(uid), time.Now())
				amount.ChatAmount += welfare
				tx.Save(&amount)

				tx.Create(&model.AccountRecord{
					Uid:           uint32(uid),
					RecordId:      uint32(t),
					Way:           1,
					Type:          "welfare",
					Amount:        welfare,
					CurrentAmount: amount.ChatAmount - amount.ChatUse,
				})

			}
		}
	}

	totalAmount := amount.ChatAmount - amount.ChatUse
	tx.Commit()

	return &types.TaskCompleteResponse{
		Total:   totalAmount,
		Amount:  add,
		Welfare: welfare,
	}, nil
}
