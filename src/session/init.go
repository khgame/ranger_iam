package session

import (
	"net/http"

	"github.com/bagaking/goulp/wlog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/khgame/ranger_iam/pkg/auth"
	"github.com/khgame/ranger_iam/src/model"
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
	group.GET("/validate", svr.HandleValidate)
}

// HandleValidate 处理验证请求
// @Summary 会话验证接口
// @Description 验证用户的JWT是否有效；根据策略选择长短票；默认长票，降级时下发短票指令
// @Tags session
// @Accept  json
// @Produce  json
// @Param Authorization header string true "带有Bearer的Token"
// @Success 200 {object} map[string]interface{} "验证成功返回用户UID"
// @Failure 401 "无效或过期的Token"
// @Router /session/validate [get]
func (svr *Service) HandleValidate(c *gin.Context) {
	token, err := auth.GetTokenStrFromHeader(c)
	if err != nil {
		wlog.ByCtx(c).WithError(err).Error("get token from header failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := svr.JWT.ValidateClaims(token)
	if err != nil {
		wlog.ByCtx(c).WithError(err).Error("validate token failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 返回创建成功的用户信息（注意不返回密码等敏感信息）
	c.JSON(http.StatusOK, gin.H{
		"uid": claims.UID,
	})
}
