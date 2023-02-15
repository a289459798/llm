// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	chat "chatgpt-tools/internal/handler/chat"
	code "chatgpt-tools/internal/handler/code"
	creation "chatgpt-tools/internal/handler/creation"
	divination "chatgpt-tools/internal/handler/divination"
	game "chatgpt-tools/internal/handler/game"
	image "chatgpt-tools/internal/handler/image"
	report "chatgpt-tools/internal/handler/report"
	user "chatgpt-tools/internal/handler/user"
	"chatgpt-tools/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/users/login",
				Handler: user.LoginHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/users/:id",
				Handler: user.UserInfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthAndUse},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/report/day",
					Handler: report.DayHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/report/week",
					Handler: report.WeekHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/report/plot",
					Handler: report.PlotHandler(serverCtx),
				},
			}...,
		),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthAndUse},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/images/create",
					Handler: image.CreateHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/images/watermark",
					Handler: image.WatermarkHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthAndUse},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/divination/qiming",
					Handler: divination.QimingHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/divination/jiemeng",
					Handler: divination.JiemengHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/divination/suanming",
					Handler: divination.SuanmingHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/divination/gsqiming",
					Handler: divination.GsqimingHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/divination/yyqiming",
					Handler: divination.YyqimingHandler(serverCtx),
				},
			}...,
		),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthAndUse},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/game/idiom",
					Handler: game.IdiomHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/game/idiom/answer",
					Handler: game.IdiomAnswerHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/game/twenty-four",
					Handler: game.TwentyFourHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/game/twenty-four/answer",
					Handler: game.TwentyFourAnswerHandler(serverCtx),
				},
			}...,
		),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthAndUse},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/code/regular",
					Handler: code.RegularHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/code/generate",
					Handler: code.GenerateHandler(serverCtx),
				},
			}...,
		),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthAndUse},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/creation/activity",
					Handler: creation.ActivityHandler(serverCtx),
				},
			}...,
		),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.AuthAndUse},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/chat/introduce",
					Handler: chat.IntroduceHandler(serverCtx),
				},
			}...,
		),
	)
}
