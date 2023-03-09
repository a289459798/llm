package task

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
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
	var openAmount int64
	var total int
	for _, m := range result {
		t, _ := strconv.Atoi(m["total"].(string))
		total += t
		if m["type"] == "share" {
			shareCount = m["count"].(int64)
		} else if m["type"] == "ad" {
			adCount = m["count"].(int64)
		} else if m["type"] == "open" {
			openAmount = m["count"].(int64)
		}
	}

	// 获取订阅
	var subscribeCount int64
	l.svcCtx.Db.
		Where("uid = ?", uid).
		Where("way = 1").
		Where("type = ?", "subscribe").
		Count(&subscribeCount)

	return &types.TaskResponse{
		Max:     10000 + 30 + 55 + 10,
		Have:    total,
		Tips:    "限时福利：今日完成所有任务额外可得10000次",
		Content: "分享每次可获得5次，被分享人第一次打开可额外获得5次，\n分享任务每天可做3次，通过分享最多可获得30次\n\n完整看完激励视频第一次获得10次，之后每个视频奖励5次\n，每天最多完成10次\n\n打开小程序即可获得10次，连续登录额外奖励对应的连续\n次数，坚持使用100年，当天可获得36500次\n\n以上所得次数仅限当天有效",
		List: []types.Task{
			{
				Title:          "分享",
				Status:         shareCount == 3,
				Total:          3,
				CompleteNumber: int(shareCount),
				Type:           "share",
				Amount:         5,
			},
			{
				Title:          "看广告",
				Status:         adCount == 10,
				Total:          10,
				CompleteNumber: int(adCount),
				Type:           "ad",
				Amount: func(count int) int {
					if count > 0 {
						return 5
					}
					return 10
				}(int(adCount)),
			},
			{
				Title:          "打开小程序",
				Status:         openAmount > 0,
				Total:          1,
				CompleteNumber: 1,
				Type:           "open",
				Amount:         int(openAmount),
			},
		},
	}, nil
}
