package auth

import (
	"errors"
	"time"

	"github.com/khicago/irr"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwtCustomClaims 包含JWT的声明
type JwtCustomClaims struct {
	UID uint `json:"uid"`
	jwt.StandardClaims
}

// JWTService 提供JWT令牌的服务
type JWTService struct {
	secretKey string
	issuer    string
}

// NewJWTService 创建JWT服务的新实例
func NewJWTService(secretKey, issuer string) *JWTService {
	return &JWTService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

// GenerateToken 生成JWT令牌
func (s *JWTService) GenerateToken(userID uint) (string, error) {
	// 设置JWT声明
	claims := &JwtCustomClaims{
		userID, // 用户ID从数据库用户模型中带入
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // 举例：让令牌在72小时后过期
			Issuer:    s.issuer,
		},
	}

	// 使用HMAC SHA256算法进行令牌签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}

// ValidateToken performs local token validation using the jwtService
func (s *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 在这里验证token方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
	// 可能存在解析错误或令牌无效错误
	if err != nil {
		return nil, irr.Wrap(err, "token claims parse failed")
	}

	// 返回验证通过的令牌
	if _, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return token, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

// ValidateClaims performs local token validation using the jwtService, return claims
func (s *JWTService) ValidateClaims(tokenStr string) (*JwtCustomClaims, error) {
	token, err := s.ValidateToken(tokenStr)
	if err != nil {
		return nil, irr.Wrap(err, "token validate failed")
	}

	return token.Claims.(*JwtCustomClaims), nil
}
