package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khicago/irr"
)

const (
	BearerSchema        = "Bearer "
	HeaderAuthorization = "Authorization"
)

// GetTokenStrFromHeader 从 header 里获取 token 字符串
func GetTokenStrFromHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(HeaderAuthorization)
	if header == "" {
		return "", irr.Error("authorization header is required")
	}

	tokenStr := strings.TrimPrefix(header, BearerSchema)
	if tokenStr == header {
		return "", irr.Error("authorization schema is wrong")
	}
	return tokenStr, nil
}

// SetTokenStrToHeader 设置 header 里的 token 字符串
func SetTokenStrToHeader(req *http.Request, token string) {
	req.Header.Set(HeaderAuthorization, BearerSchema+token)
}
