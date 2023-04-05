// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	callbackpay "chatgpt-tools/internal/handler/callback/pay"
	chat "chatgpt-tools/internal/handler/chat"
	chatbrain "chatgpt-tools/internal/handler/chat/brain"
	code "chatgpt-tools/internal/handler/code"
	common "chatgpt-tools/internal/handler/common"
	convert "chatgpt-tools/internal/handler/convert"
	creation "chatgpt-tools/internal/handler/creation"
	crontab "chatgpt-tools/internal/handler/crontab"
	divination "chatgpt-tools/internal/handler/divination"
	game "chatgpt-tools/internal/handler/game"
	hashrate "chatgpt-tools/internal/handler/hashrate"
	image "chatgpt-tools/internal/handler/image"
	order "chatgpt-tools/internal/handler/order"
	report "chatgpt-tools/internal/handler/report"
	user "chatgpt-tools/internal/handler/user"
	userai "chatgpt-tools/internal/handler/user/ai"
	userhistory "chatgpt-tools/internal/handler/user/history"
	usernotify "chatgpt-tools/internal/handler/user/notify"
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
			{
				Method:  http.MethodGet,
				Path:    "/users/login/wxqrcode",
				Handler: user.LoginWXQrcodeHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/users/login/check",
				Handler: user.LoginCheckHandler(serverCtx),
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
					Method:  http.MethodDelete,
					Path:    "/users/history/chat",
					Handler: userhistory.CleanChatListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/users/history/suanli",
					Handler: userhistory.SuanliListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/users/history/tools",
					Handler: userhistory.ToolsListHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/users/history/hashrate-exchange",
					Handler: userhistory.ExchangeHashRateListHandler(serverCtx),
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
					Path:    "/users/notify/unread",
					Handler: usernotify.UnreadHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/users/notify",
					Handler: usernotify.ListHandler(serverCtx),
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
				{
					Method:  http.MethodPost,
					Path:    "/images/pic-repair",
					Handler: image.Old2newHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/images/img2text",
					Handler: image.Image2TextHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/images/pic-repair-async",
					Handler: image.Old2newAsyncHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/images/task",
				Handler: image.TaskHandler(serverCtx),
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
			{
				Method:  http.MethodGet,
				Path:    "/common/shortlink",
				Handler: common.ShortLinkHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/common/upload-token",
				Handler: common.UploadTokenHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/common/valid/image",
				Handler: common.ValidImageHandler(serverCtx),
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
				{
					Method:  http.MethodGet,
					Path:    "/vip/give",
					Handler: vip.VipGiveHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/vip/exchange",
					Handler: vip.VipCxchangeHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/vip/code-generate",
					Handler: vip.VipCodeGenerateHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/vip/privilege",
					Handler: vip.VipPrivilegeHandler(serverCtx),
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
					Path:    "/hashrate/price",
					Handler: hashrate.HashRatePriceHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/hashrate/code-generate",
					Handler: hashrate.HashRateCodeGenerateHandler(serverCtx),
				},
				{
					Method:  http.MethodPost,
					Path:    "/hashrate/exchange",
					Handler: hashrate.HashRateCxchangeHandler(serverCtx),
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
				Method:  http.MethodGet,
				Path:    "/wechat/event/:appkey",
				Handler: wechat.ValidateHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/wechat/event/:appkey",
				Handler: wechat.EventHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/wechat/callback/subscribe",
				Handler: wechat.SubscribeCallHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/wechat/menu",
				Handler: wechat.SetMenuHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/callback/pay/vip/:type/:merchant",
				Handler: callbackpay.PayVipHandler(serverCtx),
			},
		},
	)

	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CronMiddle},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/crontab/vip-check",
					Handler: crontab.VipCheckHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/crontab/hashrate-check",
					Handler: crontab.HashRateCheckHandler(serverCtx),
				},
			}...,
		),
	)
}
