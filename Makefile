createdb:
	docker exec -it postgres12 createdb -U root --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
gen_mock:
	mockgen -package mockDb -destination db/mock/store.go github.com/mikeheft/go-backend/db/sqlc Store
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:54321/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:54321/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:54321/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:54321/simple_bank?sslmode=disable" -verbose down 1
postgres:
	docker run --name postgres12 --network bank-network -p 54321:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret/password -d postgres:12-alpine
server:
	go run main.go
sqlc:
	docker run --rm -v $(PWD):/src -w /src sqlc/sqlc generate
test:
	go test -v -cover ./...


.PHONY: postgres createdb dropdb migrateup migratedown server sqlc test migrateup1 migratedown1