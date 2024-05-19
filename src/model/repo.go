package model

import (
	"context"

	"github.com/khgame/ranger_iam/internal/util"
	"github.com/khicago/irr"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

// Repo 提供用户注册的服务
type Repo struct {
	DB *gorm.DB
}

// NewRepo 创建一个新的RegisterService实例
func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

// RegisterParams 定义注册参数结构体
type RegisterParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register 创建一个新用户
func (repo *Repo) Register(ctx context.Context, params RegisterParams) (*User, error) {
	// 通常你还需要在这里加密密码，这里为了简化示例，我们略过这一步
	// passwordHash := HashPassword(params.Password)
	id, err := util.GenIDU64(ctx)
	if err != nil {
		return nil, irr.Wrap(err, "generate id for new user failed")
	}

	user := &User{
		I64Model: I64Model{
			ID: id,
		},
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: params.Password, // 使用passwordHash代替明文密码
	}

	// 使用Gorm创建新的用户记录
	user, err = repo.createUser(ctx, user, nil, nil)
	if err != nil {
		return nil, irr.Wrap(err, "create user failed")
	}

	return user, nil
}

// FindUserByName 用过名字找到一个新用户
func (repo *Repo) FindUserByName(username string) (*User, error) {
	var user User
	err := repo.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *Repo) createUser(ctx context.Context, user *User, oauth *OAuthCredential, twoFactor *TwoFactorSetting) (*User, error) {
	// 开始数据库事务
	tx := repo.DB.WithContext(ctx).Begin()

	if twoFactor != nil {
		twoFactor.ID = user.ID
		// 首先创建TwoFactorSetting，获取生成的ID
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(twoFactor).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 接下来创建用户本身
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if oauth != nil {
		// 然后设置OAuthCredential的UserID外键
		oauth.ID = user.ID

		// 最后创建OAuthCredential条目
		if err := tx.Create(&oauth).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 所有操作成功后提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 返回创建的用户对象
	return user, nil
}
