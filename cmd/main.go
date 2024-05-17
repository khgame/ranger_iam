// File: cmd/main.go

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"github.com/khgame/ranger_iam/internal/util"
	"github.com/khgame/ranger_iam/src/app"
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

	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(dsn()), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 其他程序初始化逻辑...

	// 初始化HTTP路由
	router := gin.Default()

	// 注入db实例到注册处理函数中
	group := router.Group("/api/v1")
	app.RegisterRoutes(group, db) // 注意: RegisterRoutes 函数签名需要接受 *gorm.DB 参数

	// 开启HTTP服务
	if err = router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
