package authcli

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khicago/irr"

	"github.com/khgame/ranger_iam/pkg/auth"
)

const UserCtxKey = "UserID"

var (
	ErrValidateRemoteStatusFailed = irr.Error("validate remote status failed")
	ErrValidateRemoteDegraded     = irr.Error("validate remote degraded")
)

// Cli represents the client used to interact with the IAM server and perform local JWT verification
type Cli struct {
	localJWT    *auth.JWTService
	AuthNSvrURL string

	httpClient *http.Client
}

// New creates a new instance of the IAM client
func New(secretKey, sessionServiceURL string) *Cli {
	jwtService := auth.NewJWTService(secretKey, "UNKNOWN")
	return &Cli{
		localJWT:    jwtService,
		AuthNSvrURL: sessionServiceURL,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// GinMW 创建一个检查 JWT 是否有效的 Gin 中间件
func (cli *Cli) GinMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 tokenStr
		tokenStr, err := auth.GetTokenStrFromHeader(c)
		if err != nil {
			// todo: 从 Cookie 中获取 tokenStr
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}

		uid, err := cli.AuthN(c, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
			return
		}
		c.Set(UserCtxKey, uid)
		c.Next()
	}
}

func GetUID(c *gin.Context) (userID uint, exists bool) {
	val, exists := c.Get(UserCtxKey)
	return val.(uint), exists
}
