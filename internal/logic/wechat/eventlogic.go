package wechat

import (
	"context"
	"fmt"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EventLogic {
	return &EventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EventLogic) Event(req types.WechatValidateRequest, r *http.Request, w http.ResponseWriter) (resp *types.WeChatCallbackResponse, err error) {

	fmt.Println(req.AppKey)
	fmt.Println("事件")

	return
}
