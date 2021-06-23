SCHEMA=./migrations
DB='postgres://postgres:123matan123@127.0.0.1:5432/sellerx_task?sslmode=disable'

migrate_up:
	migrate -path $(SCHEMA) -database $(DB) up

migrate_down:
	migrate -path $(SCHEMA) -database $(DB) down