# Define the docker-compose command
DC := cd dev_ranger_iam && docker-compose

# Services names
MYSQL_SERVICE := mysql
REDIS_SERVICE := redis
MIGRATION_SERVICE := migration

# Paths
MIGRATE_UP_PATH := migration/migrate_up.sql
MIGRATE_DOWN_PATH := migration/migrate_down.sql

.PHONY: default help gen-doc bundle compose-up compose-down compose-re compose-logs db-migrate-up db-migrate-down db-migrate-re dev build-dev build-app dev-logs dev-shell

# Default to help
default: help

# Show help
help:
	@echo "Available commands:"
	@echo "  make compose-up : Start all services with docker-compose"
	@echo "  make compose-down : Stop all services and remove containers"
	@echo "  make compose-logs : Fetch logs for all services"
	@echo "  make db-migrate-up : Perform database migrations"
	@echo "  make db-migrate-down : Rollback database migrations"
	@echo "  make compose-reset : Stop all services and remove data"

gen-doc:
	@echo '```' > ./doc/PROJECT_STRUCTURE.md
	@tree -I 'bundle*' --dirsfirst --noreport >> ./doc/PROJECT_STRUCTURE.md
	@echo '```' >> ./doc/PROJECT_STRUCTURE.md
	swag init -g cmd/main.go -o doc/

# bundle, @see github.com/bagaking/file_bundle
bundle: gen-doc
	$(MAKE) -C bundle -f Makefile clean
	$(MAKE) -f bundle/Makefile
	#file_bundle -v -i ./bundle/_.file_bundle_rc -o ./bundle/_.bundle.txt

# Start services
compose-up:
	$(DC) up -d

# Stop services
compose-down:
	$(DC) down

compose-re: compose-down compose-up

# Display logs
compose-logs:
	$(DC) logs

# Existing makefile content
define run-migration
    $(DC) exec $(MYSQL_SERVICE) bash /migrate.sh $(1)
endef
#$(DC) run -T --rm migrator
#$(DC) run --rm -e MIGRATE_DIRECTION=down migrator

db-migrate-up:
	$(call run-migration,up)

db-migrate-down:
	$(call run-migration,down)

db-migrate-re: db-migrate-down db-migrate-up

# Clean up environment including volumes
compose-reset:
	$(DC) down --volumes

# Build the docker image for our go application
build-app: bundle
	$(DC) build app

dev:
	$(DC) up -d app
	@echo Starting app service...
	#$(DC) exec -d app /app/ranger_iam
	@echo App service has been started in the background.

# Start the built go application
build-dev: gen-doc build-app compose-re dev

# Connect to the app container and follow the logs
dev-logs:
	tail -f ./dev_ranger_iam/logs/ranger_iam.log

dev-shell:
	$(DC) exec app /bin/sh