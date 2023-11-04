.PHONY: postgres adminer migrate

postgres:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=pass -d postgres

adminer:
	docker run --rm -ti --network host adminer

migrate:
	migrate -source file://migrations -database postgres://postgres:pass@127.0.0.1:5432?sslmode=disable up

migrate-down:
	migrate -source file://migrations -database postgres://postgres:pass@127.0.0.1:5432?sslmode=disable down