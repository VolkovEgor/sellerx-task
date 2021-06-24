APP=cmd/main.go

build:
	go build -o bin/app.out $(APP)

run:
	go run $(APP)

lint:
	go fmt ./...
	golangci-lint run

swag:
	swag init --parseDependency -d ./internal/delivery -o ./docs/swagger -g handler.go

SCHEMA=./migrations
DB='postgres://postgres:123matan123@127.0.0.1:5432/sellerx_task?sslmode=disable'

migrate_up:
	migrate -path $(SCHEMA) -database $(DB) up

migrate_down:
	migrate -path $(SCHEMA) -database $(DB) down