.PHONY: run build test lint migrate coverage

run:
	go run ./cmd/server/

build:
	go build -o bin/server ./cmd/server/

test:
	go test ./...

lint:
	golangci-lint run

migrate:
	@echo "Migrations run automatically via GORM AutoMigrate on server start."
	@echo "Run 'make run' to apply migrations."

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out
