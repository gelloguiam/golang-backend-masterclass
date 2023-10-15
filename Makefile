postgres-run:
	docker run --name postgres12 -p 5499:5432 -e POSTGRES_USERNAME=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine

postgres-start:
	docker start postgres12

postgres-rm:
	docker rm postgres12

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

sqlc: 
	sqlc generate

.PHONY: postgres-run, postgres-start, postgres-rm, createdb, dropdb, sqlc