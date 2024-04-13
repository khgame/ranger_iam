// File: src/app/gw/routes.go

package app

import (
	"github.com/khgame/ranger_iam/internal/util"
	"github.com/khgame/ranger_iam/pkg/auth"
	"github.com/khgame/ranger_iam/src/passport"
	"github.com/khgame/ranger_iam/src/session"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes - routers all in one
// todo: using rpc
func RegisterRoutes(router gin.IRouter, db *gorm.DB) {

	// todo: 这些值应该从配置中安全获取，现在 MVP 一下
	jwtService := auth.NewJWTService("my_secret_key", util.DefaultJWTIssuer)
	//nwAuth := jwtService.GinMW()

	authGroup := router.Group("/auth")
	{
		svrPassport, _ := passport.Init(db, jwtService)
		svrPassport.ApplyMux(authGroup)
	}

	sessionGroup := router.Group("/session")
	{
		svrPassport, _ := session.Init(db, jwtService)
		svrPassport.ApplyMux(sessionGroup)
	}

}
