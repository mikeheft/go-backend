createdb:
	docker exec -it postgres12 createdb -U root --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:54321/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:54321/simple_bank?sslmode=disable" -verbose down
postgres:
	docker run --name postgres12 -p 54321:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
sqlc:
	docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate


.PHONY: postgres createdb dropdb migrateup migratedown