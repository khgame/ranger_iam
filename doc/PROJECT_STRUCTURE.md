```
.
├── cmd
│   └── main.go
├── config
│   ├── app.dev.yaml
│   ├── app.prod.yaml
│   ├── log.dev.yaml
│   └── log.prod.yaml
├── deployment
│   ├── migration
│   │   ├── 01.passport_down.sql
│   │   ├── 01.passport_up.sql
│   │   └── migrate.sh
│   ├── Dockerfile
│   └── ci_cd.yaml
├── dev_ranger_iam
│   ├── dial
│   │   ├── auth.http
│   │   └── http-client.env.json
│   └── docker-compose.yml
├── doc
│   └── PROJECT_STRUCTURE.md
├── internal
│   ├── repository
│   │   ├── cache.go
│   │   ├── dc.go
│   │   └── rds.go
│   └── util
│       ├── const.go
│       └── env.go
├── pkg
│   ├── auth
│   │   ├── jwt.go
│   │   ├── oauth.go
│   │   └── util.go
│   └── authcli
│       ├── cli.go
│       ├── refresh.go
│       └── validate.go
├── script
│   └── setup_project.sh
├── src
│   ├── app
│   │   ├── error_handler.go
│   │   └── routes.go
│   ├── model
│   │   ├── repo.go
│   │   └── user.go
│   ├── passport
│   │   ├── init.go
│   │   ├── login.go
│   │   ├── register.go
│   │   └── util.go
│   └── session
│       ├── init.go
│       ├── longterm.go
│       └── shortterm.go
├── LICENSE
├── MODULES.puml
├── Makefile
├── README.md
├── go.mod
└── go.sum
```
