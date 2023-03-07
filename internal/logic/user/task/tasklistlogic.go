package task

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
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
		Select("type, count(*) as count").
		Where("uid = ?", uid).
		Where("created_at between ? and ?", today+" 00:00:00", today+" 23:59:59").
		Where("way = 1").
		Group("type").
		Scan(&result)

	var shareCount int64
	var adCount int64
	var openAmount int64
	for _, m := range result {
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
		Content: "分享每次可获得5次，通过分享链接打开可额外获得5次，\n分享任务每天可做3次，通过分享最多可获得30次\n\n完整看完激励视频可获得10次\n\n关注公众号可获得20次",
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
				Amount:         10,
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
