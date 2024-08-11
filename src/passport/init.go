package passport

import (
	"github.com/gin-gonic/gin"
	"github.com/khgame/ranger_iam/pkg/auth"
	"github.com/khgame/ranger_iam/src/model"
	"gorm.io/gorm"
)

type Service struct {
	Repo *model.Repo
	JWT  *auth.JWTService
}

var svr *Service

func Init(db *gorm.DB, jwtService *auth.JWTService) (*Service, error) {
	svr = &Service{
		Repo: model.NewRepo(db),
		JWT:  jwtService,
	}
	return svr, nil
}

func (svr *Service) ApplyMux(group gin.IRouter) {
	group.POST("/register", svr.HandleRegister)
	group.POST("/login", svr.HandleLogin)
	group.POST("/wechat/miniapp", svr.HandleLoginWechatMiniApp)
}
