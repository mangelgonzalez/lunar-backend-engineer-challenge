# Go parameters
MAIN_PATH=cmd/rockets/

.PHONY: setup
setup: SHELL:=/bin/bash
setup: ## Setup Project
setup:
	@if [ -f ./docker-compose.yaml ]; then \
		echo "Detected 'docker-compose.yaml' file, moving it to 'docker-compose.yaml.old'"; \
		mv ./docker-compose.yaml ./docker-compose.yaml.old && \
		cp docker-compose.yaml.dist docker-compose.yaml; \
	else \
		cp docker-compose.yaml.dist docker-compose.yaml; \
	fi
	@if [ -f ./.env ]; then \
		echo "Detected '.env' file, moving it to '.env.old'"; \
		mv ./.env ./.env.old && \
		cp .env.dist .env; \
	else \
		cp .env.dist .env; \
	fi

	docker-compose pull
	docker-compose build

.PHONY: start
start: ## Start the local environment
	[ -f ./docker-compose.yaml ] && [ -f ./.env ] || $(MAKE) setup
	docker-compose up -d

.PHONY: stop
stop: ## Stop the local environment
	docker-compose down

.PHONY: test
test: unit-tests acceptance-tests

.PHONY: unit-tests
unit-tests: ## Run unitary tests
	go test ./internal/... -race

.PHONY: acceptance-tests
acceptance-tests: ## Run acceptance tests
	go test ./test/... -race

.PHONY: migrateup
migrateup: ## Run migration up
	go run cmd/rockets/database-migrations/migrations.go "up"

.PHONY: migratedown
migratedown: ## Run migration down
	go run cmd/rockets/database-migrations/migrations.go "down"

.PHONY: newmigration
newmigration:
	sql-migrate new -config test/migrations/dbconfig.yml

help: ## Display this help screen
	@echo "Usage:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
