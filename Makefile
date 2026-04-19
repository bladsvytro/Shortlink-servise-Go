.PHONY: help build run test clean migrate seed lint docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  build     - Build the application"
	@echo "  run       - Run the application"
	@echo "  test      - Run tests"
	@echo "  lint      - Run linter"
	@echo "  clean     - Clean build artifacts"
	@echo "  migrate   - Run database migrations"
	@echo "  seed      - Seed database with test data"
	@echo "  docker-up - Start Docker containers"
	@echo "  docker-down - Stop Docker containers"

build:
	go build -o bin/url-shortener ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	golangci-lint run

clean:
	rm -rf bin/ coverage.out

migrate:
	@echo "Running migrations..."
	# TODO: Implement migration command
	# migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

seed:
	@echo "Seeding database..."
	# TODO: Implement seed command

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker build -t url-shortener:latest .

docker-run:
	docker run -p 8080:8080 --env-file .env url-shortener:latest

dev:
	air -c .air.toml

generate:
	@echo "Generating code..."
	# TODO: Add code generation commands

check: lint test

.PHONY: all
all: check build