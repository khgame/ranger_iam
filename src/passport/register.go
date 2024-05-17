package passport

import (
	"github.com/gin-gonic/gin"
	"github.com/khgame/ranger_iam/src/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// RegisterRequest 定义注册请求的结构
type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// ErrorMessage 定义错误信息的结构
type ErrorMessage struct {
	Message string `json:"message"`
}

// HandleRegister 处理注册请求
// @Summary 用户注册接口
// @Description 用户填入用户名、邮箱和密码进行注册
// @Tags auth
// @Accept  json
// @Produce  json
// @Param register body RegisterRequest true "注册请求信息"
// @Success 201 {object} map[string]any "注册成功返回新创建的用户信息和token"
// @Failure 400 "请求格式错误或密码不匹配"
// @Failure 500 "无法注册用户或生成token"
// @Router /auth/register [post]
func (svr *Service) HandleRegister(c *gin.Context) {
	var req RegisterRequest

	// 绑定请求体到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// 验证两次输入的密码是否匹配
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not encrypt password"})
		return
	}

	// 实例化注册服务并进行注册
	user, err := svr.Repo.Register(c, model.RegisterParams{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		// 处理可能的数据库错误，如唯一性违反等
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user, " + err.Error()})
		return
	}

	// 注册成功后生成JWT (short-ticket sample)
	token, err := svr.genJWTTokenAndSetCookie(c.Writer, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token, " + err.Error()})
		return
	}

	// 返回创建成功的用户信息（注意不返回密码等敏感信息）
	c.JSON(http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}
