DB_URL=postgres://postgres:12345@localhost:5432/aulway?sslmode=disable

migrate-up:
	migrate -database $(DB_URL) -path internal/database/postgres/migration up

migrate-down:
	migrate -database $(DB_URL) -path internal/database/postgres/migration down

create-migration:
	@read -p "migration name: " name; \
	migrate create -ext sql -dir internal/database/postgres/migration -tz "UTC" $$name

initswag:
	swag init --parseDependency --parseInternal --propertyStrategy pascalcase --parseDepth 3 -g main.go

stripe:
	brew install stripe/stripe-cli/stripe # Mac
	docker run -d -p 12111:12111 stripe/stripe-mock # Docker

