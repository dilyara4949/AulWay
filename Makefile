DB_URL=postgres://postgres:12345@localhost:5432/aulway?sslmode=disable

migrate-up:
	migrate -database $(DB_URL) -path internal/database/postgres/migration up

migrate-one-down:
	migrate -database $(DB_URL) -path internal/database/postgres/migration down 1

migrate-down:
	migrate -database $(DB_URL) -path internal/database/postgres/migration down

create-migration:
	@read -p "migration name: " name; \
	migrate create -ext sql -dir internal/database/postgres/migration -tz "UTC" $$name

initswag:
	swag init --parseDependency --parseInternal --propertyStrategy pascalcase --parseDepth 3 -g main.go

net:
	docker network create my_network docker run --name postgre --network my_network -e POSTGRES_PASSWORD=12345 -p 5435:5435 -d postgres

compose:
	docker-compose -f docker/docker-compose.yml up  -d