package model

import (
	"time"

	"gorm.io/gorm"
)

type I64Model struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// User 定义了用户账号管理模块的模型
type User struct {
	I64Model

	Username     string `gorm:"uniqueIndex;not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
}

type OAuthCredential struct {
	I64Model

	Provider   string `gorm:"not null"` // local, google, facebook, twitter, wechat
	ProviderID string `gorm:"index"`
}

type TwoFactorSetting struct {
	I64Model

	IsEnabled      bool
	Phone          string
	SecondaryEmail string
	BackupCodes    string // JSON array of backup codes
}

// BeforeCreate 是Gorm的hook，在创建记录前调用
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate 是Gorm的hook，在更新记录前调用
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}
