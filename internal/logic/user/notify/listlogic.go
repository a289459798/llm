package notify

import (
	"chatgpt-tools/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List() (resp *types.UserNotifyListResponse, err error) {
	uid, _ := l.ctx.Value("uid").(json.Number).Int64()
	list := []model.AINotify{}
	l.svcCtx.Db.Where("uid = ?", uid).Limit(20).Order("status asc, id desc").Find(&list)
	if len(list) == 0 {
		return nil, errors.New(string(http.StatusNotFound))
	}

	res := []types.UserNotifyResponse{}
	for _, notify := range list {
		res = append(res, types.UserNotifyResponse{
			Title:    notify.Title,
			Contennt: notify.Content,
			Status:   notify.Status,
		})
	}

	l.svcCtx.Db.Model(&model.AINotify{}).Where("uid = ?", uid).Update("status", 1)

	return &types.UserNotifyListResponse{Data: res}, nil
}
