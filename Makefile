APP=cmd/main.go

build:
	go build -o bin/app.out $(APP)

run:
	go run $(APP) local_config

SCHEMA=./migrations
DB='postgres://postgres:1234@127.0.0.1:5432/sellerx_task?sslmode=disable'

migrate_up:
	migrate -path $(SCHEMA) -database $(DB) up

migrate_down:
	migrate -path $(SCHEMA) -database $(DB) down