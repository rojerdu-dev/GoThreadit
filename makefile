.PHONY: postgres adminer migrate

postgres:
	docker run --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=pass -d postgres

adminer:
	docker run --name adminer -p 8008:8080 -d adminer

migrate:
	migrate -source file://migrations -database postgres://postgres:pass@127.0.0.1:5432?sslmode=disable up

migrate-down:
	migrate -source file://migrations -database postgres://postgres:pass@127.0.0.1:5432?sslmode=disable down