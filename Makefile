postgres-run:
	docker run --name postgres12 -p 5499:5432 -e POSTGRES_USERNAME=postgres -e POSTGRES_PASSWORD=password -d postgres:12-alpine

postgres-start:
	docker start postgres12

postgres-rm:
	docker rm postgres12

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrate:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5499/simple_bank?sslmode=disable" -verbose up

sqlc:
	sqlc generate

test:
	go test -v -count=30 -cover ./...

.PHONY: postgres-run, postgres-start, postgres-rm, createdb, dropdb, migrate, sqlc, test