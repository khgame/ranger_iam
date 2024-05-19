// File: cmd/main.go

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bagaking/goulp/wlog"
	"github.com/gin-gonic/gin"
	"github.com/khgame/ranger_iam/internal/util"
	"github.com/khgame/ranger_iam/src/app"
	"github.com/khicago/irr"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// dsn for dev
// todo: using env config
func dsn() string {
	// 使用docker-compose环境变量来设置数据库DSN
	username := "user"
	password := "password"

	// todo: using config
	host := "localhost"
	switch util.Env() {
	case util.RuntimeENVDev:
		host = "mysql"
	case util.RuntimeENVProd:
		host = "mysql"
	case util.RuntimeENVLocal:
		fallthrough
	default:
	}

	port := "3306"
	dbname := "ranger_iam"
	charset := "utf8mb4"
	loc := "Local"
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s", username, password, host, port, dbname, charset, loc)
}

func main() {
	// 配置一个lumberjack.Logger
	logRoller := &lumberjack.Logger{
		Filename:   "./logs/ranger_iam.log", // 日志文件的位置
		MaxSize:    10,                      // 日志文件的最大大小（MB）
		MaxBackups: 31,                      // 保存的旧日志文件最大个数
		MaxAge:     31,                      // 保存的旧日志文件的最大天数
		Compress:   true,                    // 是否压缩归档的日志文件
	}
	defer func() {
		if err := logRoller.Close(); err != nil {
			fmt.Println("Failed to close log", err)
		}
	}()
	mustInitLogger(logRoller)

	// 初始化数据库连接
	db := mustInitDB()
	// 其他程序初始化逻辑...

	// 初始化HTTP路由
	router := gin.Default()
	router.Use(
		ginRecoveryWithLog(),
	)

	// 注入db实例到注册处理函数中
	group := router.Group("/api/v1")
	app.RegisterRoutes(group, db) // 注意: RegisterRoutes 函数签名需要接受 *gorm.DB 参数

	wlog.Common("main").Infof("service start ..")
	// 开启HTTP服务
	if err := router.Run(":8080"); err != nil {
		wlog.Common("main").WithError(err).Infof("gin exit")
	}
}

func mustInitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn()), &gorm.Config{})
	if err != nil {
		wlog.Common("main", "mustInitDB").Fatal("failed to connect database:", err)
	}
	return db
}

func mustInitLogger(fileLogger io.Writer) {
	multiLogger := io.MultiWriter(fileLogger, os.Stdout)
	logrus.SetOutput(multiLogger)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel) // 设置日志记录级别

	// Gin 设置
	gin.DisableConsoleColor()
	gin.DefaultWriter = logrus.StandardLogger().Out

	wlog.SetEntryGetter(func(ctx context.Context) *logrus.Entry {
		entry := logrus.NewEntry(logrus.StandardLogger())
		return entry.WithContext(ctx)
	})
}

// ginRecoveryWithLog 返回一个中间件，当程序发生 panic 时记录错误日志，并返回 HTTP 500 错误。
func ginRecoveryWithLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var er irr.IRR
				if e, ok := err.(error); ok {
					er = irr.TrackSkip(1, e, "recover from panic!!")
				} else {
					er = irr.TrackSkip(1, irr.Error("%v", err), "recover from panic!!")
				}
				er = er.LogError(wlog.ByCtx(c, c.HandlerName()))

				// 返回 HTTP 500 错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "Internal Server Error",
				})
			}
		}()

		// 处理请求
		c.Next()
	}
}
