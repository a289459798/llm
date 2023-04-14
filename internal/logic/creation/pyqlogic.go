package creation

import (
	"context"
	"net/http"

	"chatgpt-tools/internal/svc"
	"chatgpt-tools/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PyqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPyqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PyqLogic {
	return &PyqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PyqLogic) Pyq(req *types.PyqRequest, w http.ResponseWriter) (resp *types.CreationResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
