package sanmuai

import (
	"chatgpt-tools/internal/svc"
	"context"
)

type SanmuAI interface {
	CreateImage(req ImageCreate) (stream []string, err error)
}

type SanmuData struct {
	Ctx    context.Context
	SvcCtx *svc.ServiceContext
}

func GetAI(model string, data SanmuData) SanmuAI {
	if model == "journey" {
		return NewJourney(data.Ctx, data.SvcCtx)
	} else {
		return NewOpenAi(data.Ctx, data.SvcCtx)
	}
}
