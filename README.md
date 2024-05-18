# Ranger IAM

Ranger IAM 是一个快速搭建身份认证和管理（IAM）服务的 Go 应用。

# 快速开始

## 环境要求
- Docker
- Docker Compose
- 
## 步骤

克隆仓库到本地：
`bash git clone https://github.com/khgame/ranger_iam.git`

构建并启动服务：
`bash make build-dev`

此命令将构建 Docker 镜像并启动服务，默认监听端口为 18080。

其他命令见 Makefile

# API

Ranger IAM 提供了一系列 REST API 接口用于用户注册、登陆和会话验证等。

具体的 API 文档可以访问 doc 目录，或运行后访问 swagger 页面 `/api/v1/swagger/index.html`

# SDK

位于 `pkg/authcli`