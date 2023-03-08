package svc

import (
	"chatgpt-tools/internal/config"
	"chatgpt-tools/internal/middleware"
	"chatgpt-tools/model"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type ServiceContext struct {
	Config     config.Config
	Db         *gorm.DB
	AuthAndUse rest.Middleware
	Sign       rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, _ := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "gpt_",
			SingularTable: true,
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Warn, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		),
	})

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Account{})
	//db.AutoMigrate(&model.Record{})
	//db.AutoMigrate(&model.Feedback{})
	//db.AutoMigrate(&model.Apikey{})
	//db.AutoMigrate(&model.Pic2Pic{})
	db.AutoMigrate(&model.AccountRecord{})
	db.AutoMigrate(&model.ShareRecord{})
	db.AutoMigrate(&model.ChatTemplate{})
	return &ServiceContext{
		Config:     c,
		Db:         db,
		AuthAndUse: middleware.NewAuthAndUseMiddleware(c, db).Handle,
		Sign:       middleware.NewSignMiddleware(c).Handle,
	}
}
