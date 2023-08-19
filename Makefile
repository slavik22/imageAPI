DB_URL=postgresql://postgres:password@localhost:5432/images?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name postgres:15-alpines -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root images

dropdb:
	docker exec -it postgres dropdb simple_bank

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go github.com/slavik22/imageAPI/db/sqlc Store
