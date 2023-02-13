package svc

import (
	"chatgpt-tools/internal/config"
	"chatgpt-tools/internal/middleware"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	Db         *gorm.DB
	Auth       rest.Middleware
	AuthAndUse rest.Middleware
	GptClient  *gogpt.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	//db, _ := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		TablePrefix:   "gpt_",
	//		SingularTable: true,
	//	},
	//	Logger: logger.New(
	//		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
	//		logger.Config{
	//			SlowThreshold:             time.Second, // 慢 SQL 阈值
	//			LogLevel:                  logger.Info, // 日志级别
	//			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
	//			Colorful:                  false,       // 禁用彩色打印
	//		},
	//	),
	//})
	//
	//db.AutoMigrate(&model.User{})
	return &ServiceContext{
		Config: c,
		//Db:         db,
		Auth:       middleware.NewAuthMiddleware().Handle,
		AuthAndUse: middleware.NewAuthAndUseMiddleware().Handle,
		GptClient:  gogpt.NewClient(c.OpenAIKey),
	}
}
