package passport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// WXLoginResp  微信小程序登录 response
type WXLoginResp struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	Token      string `json:"token"`
	Auth       uint32 `json:"auth"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Nickname   string `json:"nickname"`
	Headimage  string `json:"headimage"`
}

// HandleLoginWechatMiniApp 处理微信登录请求
func (svr *Service) HandleLoginWechatMiniApp(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code is required"})
		return
	}

	var (
		appID     = os.Getenv("WECHAT_APP_ID")
		appSecret = os.Getenv("WECHAT_APP_SECRET")
	)

	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grrant_type=authorization_code"
	url = fmt.Sprintf(url, appID, appSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Validation failed"})
		return
	}
	defer resp.Body.Close()
	// 解析http请求中body数据到我们定义的结构体中
	wxResp := WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&wxResp); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "decoding failed"})
		return
	}
	// 判断微信接口是否返回异常
	if wxResp.ErrCode != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "request failed"})
		return
	}

	user, err := svr.Repo.UpsertUserByProvider(c, "wechat", wxResp.UnionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
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
