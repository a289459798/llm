package svc

import (
	"chatgpt-tools/internal/config"
	"chatgpt-tools/internal/middleware"
	"chatgpt-tools/model"
	log2 "github.com/sirupsen/logrus"
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
	CronMiddle rest.Middleware
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
				SlowThreshold: time.Second, // 慢 SQL 阈值
				LogLevel: func() logger.LogLevel {
					if c.Mode == "dev" {
						return logger.Info
					}
					return logger.Warn
				}(), // 日志级别
				IgnoreRecordNotFoundError: true, // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true, // 禁用彩色打印
			},
		),
	})

	//db.AutoMigrate(&model.AIUser{})
	//db.AutoMigrate(&model.Account{})
	//db.AutoMigrate(&model.Record{})
	//db.AutoMigrate(&model.Feedback{})
	//db.AutoMigrate(&model.Apikey{})
	//db.AutoMigrate(&model.Pic2Pic{})
	//db.AutoMigrate(&model.AccountRecord{})
	//db.AutoMigrate(&model.ShareRecord{})
	//db.AutoMigrate(&model.ChatTemplate{})
	//db.AutoMigrate(&model.Contraband{})
	//db.AutoMigrate(&model.AI{})
	//db.AutoMigrate(&model.Message{})
	//db.AutoMigrate(&model.Error{})
	//db.AutoMigrate(&model.Setting{})
	//db.AutoMigrate(&model.Order{})
	//db.AutoMigrate(&model.OrderPay{})
	//db.AutoMigrate(&model.OrderItem{})
	//db.AutoMigrate(&model.VipCode{})
	//db.AutoMigrate(&model.Vip{})
	//db.AutoMigrate(&model.User{})
	//db.AutoMigrate(&model.AIUserVip{})
	//db.AutoMigrate(&model.App{})
	//db.AutoMigrate(&model.RequestLog{})
	//db.AutoMigrate(&model.ScanScene{})
	//db.AutoMigrate(&model.AIHashRate{})
	//db.AutoMigrate(&model.AIHashRateCode{})
	//db.AutoMigrate(&model.AIUserHashRate{})
	//db.AutoMigrate(&model.AINotify{})
	db.AutoMigrate(&model.Distributor{})
	db.AutoMigrate(&model.DistributorLevel{})
	db.AutoMigrate(&model.DistributorRecord{})

	if c.Mode == "dev" {
		log2.SetLevel(log2.DebugLevel)
	} else {
		log2.SetLevel(log2.WarnLevel)
	}
	return &ServiceContext{
		Config:     c,
		Db:         db,
		AuthAndUse: middleware.NewAuthAndUseMiddleware(c, db).Handle,
		Sign:       middleware.NewSignMiddleware(c).Handle,
		CronMiddle: middleware.NewCronMiddleMiddleware().Handle,
	}
}
