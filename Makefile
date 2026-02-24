.PHONY: run build tidy lint

APP=./cmd/server

run:
	go run $(APP)/main.go

build:
	go build -o bin/server $(APP)/main.go

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down