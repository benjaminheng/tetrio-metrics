sqlc:
	sqlc generate

new-migration:
	migrate create -ext sql -dir ./migrations/ $(NAME)

migrate-up:
	migrate -source file://migrations/ -database sqlite3://data.db up

install:
	go install ./...
