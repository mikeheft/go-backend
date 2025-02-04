createdb:
	docker exec -it postgres12 createdb -U root --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
# GRPC client for local use
evans:
	evans --host localhost --port 9090 -r repl
mock:
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
	docker run --name postgres12 --network bank-network -p 54321:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine
server:
	go run main.go
sqlc:
	docker run --rm -v $(PWD):/src -w /src sqlc/sqlc generate
test:
	go test -v -cover ./...


.PHONY: postgres createdb dropdb evans migrateup migratedown server sqlc test migrateup1 migratedown1 proto redis mock