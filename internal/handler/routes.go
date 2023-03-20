// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	chat "chatgpt-tools/internal/handler/chat"
	chatbrain "chatgpt-tools/internal/handler/chat/brain"
	code "chatgpt-tools/internal/handler/code"
	common "chatgpt-tools/internal/handler/common"
	convert "chatgpt-tools/internal/handler/convert"
	creation "chatgpt-tools/internal/handler/creation"
	divination "chatgpt-tools/internal/handler/divination"
	game "chatgpt-tools/internal/handler/game"
	image "chatgpt-tools/internal/handler/image"
	order "chatgpt-tools/internal/handler/order"
	report "chatgpt-tools/internal/handler/report"
	user "chatgpt-tools/internal/handler/user"
	userai "chatgpt-tools/internal/handler/user/ai"
	userhistory "chatgpt-tools/internal/handler/user/history"
	usertask "chatgpt-tools/internal/handler/user/task"
	vip "chatgpt-tools/internal/handler/vip"
	wechat "chatgpt-tools/internal/handler/wechat"
	"chatgpt-tools/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/users/login",
				Handler: user.LoginHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/users",
				Handler: user.UserInfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Sign},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/users/tasks",
					Handler: usertask.TaskListHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/users/tasks",
					Handler: usertask.TaskCompleteHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/users/tasks/follow",
					Handler: usertask.TaskShareFollowHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/users/ai",
				Handler: userai.AiInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/users/ai",
				Handler: userai.AiEditHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Sign},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/users/history/chat",
					Handler: userhistory.ChatListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/users/history/suanli",
					Handler: userhistory.SuanliListHandler(serverCtx),
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
				{
					Method:  http.MethodPost,
					Path:    "/report/work",
					Handler: report.WorkHandler(serverCtx),
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
					Path:    "/images/create",
					Handler: image.CreateHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/images/create/multi",
					Handler: image.CreateMultiHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/images/watermark",
					Handler: image.WatermarkHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/images/edit",
					Handler: image.EditHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/images/pic2pic",
					Handler: image.Pic2picHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/images/pic2pic/task",
					Handler: image.Pic2pictaskHandler(serverCtx),
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
				{
					Method:  http.MethodPost,
					Path:    "/divination/horoscope",
					Handler: divination.HoroscopeHandler(serverCtx),
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
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
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
				{
					Method:  http.MethodPost,
					Path:    "/code/name",
					Handler: code.NameHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/code/exam",
					Handler: code.ExamHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/code/playground",
					Handler: code.PlaygroundHandler(serverCtx),
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
					Path:    "/creation/activity",
					Handler: creation.ActivityHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/creation/diary",
					Handler: creation.DiaryHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/creation/article",
					Handler: creation.ArticleHandler(serverCtx),
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
					Path:    "/chat/introduce",
					Handler: chat.IntroduceHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/chat/salary",
					Handler: chat.SalaryHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/chat/reject",
					Handler: chat.RejectHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/chat/pursue",
					Handler: chat.PursueHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Sign},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/chat/chat/:chatId",
					Handler: chatbrain.ChatHistoryHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/chat/chat/template",
					Handler: chatbrain.ChatTemplateHandler(serverCtx),
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
					Path:    "/chat/chat",
					Handler: chatbrain.ChatHandler(serverCtx),
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
					Path:    "/convert/translate",
					Handler: convert.TranslateHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/common/valid/text",
				Handler: common.ValidTextHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/common/valid/chat",
				Handler: common.ValidChatHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/common/message",
				Handler: common.MessageHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/common/qrcode",
				Handler: common.QrcodeHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Sign},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/vip/price",
					Handler: vip.VipPriceHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.Sign},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/order/vip",
					Handler: order.VipOrderCreateHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/wechat/callback/subscribe",
				Handler: wechat.SubscribeCallHandler(serverCtx),
			},
		},
	)
}
