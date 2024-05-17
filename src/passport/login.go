// File: src/profile/passport/login_handler.go

package passport

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 定义登录请求的结构
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// HandleLogin 处理登录请求
// @Summary 用户登录接口
// @Description 用户使用用户名和密码登录
// @Tags auth
// @Accept  json
// @Produce  json
// @Param login body LoginRequest true "登录请求信息"
// @Success 200 {object} map[string]interface{} "登录成功返回用户信息和token"
// @Failure 400 "请求格式错误"
// @Failure 401 "无效的用户名或密码"
// @Failure 500 "服务器内部错误"
// @Router /auth/login [post]
func (svr *Service) HandleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 查找用户，验证凭证
	user, err := svr.Repo.FindUserByName(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// 检查密码是否匹配
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 注册成功后生成JWT (short-ticket sample)
	token, err := svr.genJWTTokenAndSetCookie(c.Writer, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 发送令牌给用户
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
