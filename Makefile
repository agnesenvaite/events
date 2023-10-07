.PHONY: start stop restart destroy migrate lint ssh mocks docs

start: ## Start Docker containers
	docker compose up -d

stop: ## Stop Docker containers
	docker compose stop

restart: ## Restart Docker containers
	docker compose restart

destroy: ## Destroy Docker containers
	docker compose down

ssh: ## SSH into app container
	docker compose run --rm api sh

migrate-init: ## Run migrations
	go install github.com/rubenv/sql-migrate/...@latest && sql-migrate up

migrate: ## Run migrations
	sql-migrate up

new-migration: ## New migration
	@read -p "Enter Migration Name:" migration_name; \
	sql-migrate new $$migration_name

lint:
	golangci-lint run --deadline 5m

test:
	go test ./...

mocks: ## Generate mocks
	./scripts/generate_mocks.sh

docs: ## Generate OpenAPI specification
	swag init -g cmd/api/app.go --parseDependency --parseDepth 2 --outputTypes=go --generatedTime

help: ## Display available commands
	grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
