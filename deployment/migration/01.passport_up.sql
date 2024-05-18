-- Up migration scripts for [[Profile]] Module

-- Creating the users table
-- migration/01.passport_up.sql

-- 创建 `users` 表，用于管理用户账号信息
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID，自增长 -- TODO 改成发号器',
    `username` VARCHAR(255) UNIQUE NOT NULL COMMENT '用户名，唯一',
    `email` VARCHAR(255) UNIQUE NOT NULL COMMENT '用户邮箱，唯一',
    `password_hash` CHAR(60) NOT NULL COMMENT '用户密码哈希值，用于密码登录验证',

    -- `created_at`和`updated_at`会由Gorm 自动维护，无需手动插入或更新这些字段。
    `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '记录创建时间',
    `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '记录更新时间',

    -- Gorm 默认的软删除是利用`deleted_at`字段，如果该字段有值，则表示记录被软删除。
    `deleted_at` DATETIME(3) COMMENT '记录软删除时间',

    INDEX `idx_users_username` (`username`),
    INDEX `idx_users_email` (`email`),
    INDEX `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- OAuthCredentials table
CREATE TABLE `user_oauth_credentials` (
    `id` BIGINT UNSIGNED PRIMARY KEY,
    `provider` ENUM('local', 'google', 'facebook', 'twitter', 'wechat') NOT NULL DEFAULT 'local' COMMENT '账号提供者（本地或社交网络）',
    `provider_id` VARCHAR(255) COMMENT '社交账号提供者的唯一标识',

    INDEX `idx_users_provider_provider_id` (`provider`, `provider_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- TwoFactorSettings table
CREATE TABLE `user_two_factor_settings` (
    `id` BIGINT UNSIGNED PRIMARY KEY,
    `is_enabled` BOOLEAN DEFAULT FALSE COMMENT 'Flag indicating if two-factor is enabled',
    `phone` VARCHAR(255) COMMENT 'Phone number for two-factor authentication',
    `secondary_email` VARCHAR(255) COMMENT 'Secondary email for two-factor authentication',
    `backup_codes` TEXT COMMENT 'JSON array of backup codes for two-factor authentication'

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
