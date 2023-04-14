package creation

import (
	"context"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type XhsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewXhsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *XhsLogic {
	return &XhsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *XhsLogic) Xhs(req *types.XhsRequest, w http.ResponseWriter) (resp *types.CreationResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
