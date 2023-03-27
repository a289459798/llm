package task

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskListLogic {
	return &TaskListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskListLogic) TaskList(req *types.InfoRequest) (resp *types.TaskResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	//获取次数
	today := time.Now().Format("2006-01-02")
	var result []map[string]interface{}
	l.svcCtx.Db.Model(&model.AccountRecord{}).
		Select("type, count(*) as count, sum(amount) as total").
		Where("uid = ?", uid).
		Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").
		Where("way = 1").
		Group("type").
		Scan(&result)

	var shareCount int64
	var adCount int64
	var openAmount int
	var total int
	for _, m := range result {
		t, _ := strconv.Atoi(m["total"].(string))
		total += t
		if m["type"] == "share" {
			shareCount = m["count"].(int64)
		} else if m["type"] == "ad" {
			adCount = m["count"].(int64)
		} else if m["type"] == "open" {
			openAmount, _ = strconv.Atoi(m["total"].(string))
		}
	}

	// 限时福利
	tips := ""
	tasksSetting, err := model.Setting{Name: "task_welfare"}.Find(l.svcCtx.Db)
	if err == nil && len(tasksSetting) > 0 {
		startTime, err1 := time.ParseInLocation("2006-01-02 15:04:05", tasksSetting["start"].(string)+" 00:00:00", time.Local)
		endTime, err2 := time.ParseInLocation("2006-01-02 15:04:05", tasksSetting["end"].(string)+" 23:59:59", time.Local)
		if err1 == nil && err2 == nil && time.Now().Unix() >= startTime.Unix() && time.Now().Unix() < endTime.Unix() {
			tips = fmt.Sprintf("限时福利：今日完成所有任务额外可得%d算力", int(tasksSetting["amount"].(float64)))
		}
	}

	// 加群
	user := &model.AIUser{}
	l.svcCtx.Db.Where("uid = ?", uid).First(&user)
	isJoinGroup := user.IsJoinGroup(l.svcCtx.Db)

	return &types.TaskResponse{
		Max:     100 + 30 + 55 + 5,
		Have:    total,
		Tips:    tips,
		Content: "每天打开小程序即可获得5算力，连续登录额外奖励对应的连续次数算力，坚持使用100年，当天可获得36500算力\n\n分享后可获得5算力，链接被新用户打开可额外获得5算力，分享任务每天可做3次，通过分享最多可获得30算力\n\n每天第一次看激励视频获得10算力，之后每个视频奖励5算力，每天可观看10次\n\n以上所得算力仅限当天有效",
		List: []types.Task{
			{
				Title:          "分享",
				Status:         shareCount == 3,
				Total:          3,
				CompleteNumber: int(shareCount),
				Type:           "share",
				Amount:         5,
				Max:            30,
			},
			{
				Title:  "看广告",
				Status: adCount >= 10,
				Total:  10,
				CompleteNumber: func() int {
					if adCount > 10 {
						return 10
					}
					return int(adCount)
				}(),
				Type: "ad",
				Amount: func(count int) int {
					if count > 0 {
						return 5
					}
					return 10
				}(int(adCount)),
				Max: 55,
			},
			{
				Title:  "加群",
				Status: isJoinGroup,
				Total:  1,
				CompleteNumber: func() int {
					if isJoinGroup {
						return 1
					}
					return 0
				}(),
				Type:   "group",
				Amount: 10,
				Max:    10,
			},
			{
				Title:          "打开小程序",
				Status:         openAmount > 0,
				Total:          1,
				CompleteNumber: 1,
				Type:           "open",
				Amount:         openAmount,
				Max:            36500,
			},
		},
	}, nil
}
