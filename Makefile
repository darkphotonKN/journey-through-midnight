# Load the .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Construct the DB_STRING dynamically
DB_STRING=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

build:
	 @go build -o bin/game-server ./cmd/
	
run: build
	@./bin/game-server

dev:
	@air

# Run tests with verbose output and coverage
test:
	@go test ./... -cover

test-logs:
	@go test -v ./... -cover

test-game-preview:
	@go test ./internal/game/ --cover -coverprofile=coverage.out 
	@go tool cover -html=coverage.out

# Migration commands using golang-migrate
migrate-up:
	@migrate -path migrations -database "$(DB_STRING)" up

migrate-down:
	@migrate -path migrations -database "$(DB_STRING)" down 1

migrate-down-all:
	@migrate -path migrations -database "$(DB_STRING)" down

migrate-status:
	@migrate -path migrations -database "$(DB_STRING)" version

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make migrate-force VERSION=<version>"; \
		exit 1; \
	fi; \
	migrate -path migrations -database "$(DB_STRING)" force $(VERSION)

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-create NAME=<migration_name>"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir migrations -seq $(NAME)

.PHONY: run test migrate-up migrate-down migrate-status








