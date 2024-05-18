# constants
MIGRATE = migrate -database "postgres://postgres:123123@127.0.0.1:5432/go-halo-suster?sslmode=disable" -path ./db/migrations -verbose

.PHONY: default
default: help

# Start dev server.
start:
	air

.PHONY: migrate-up
migrate-up:
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down:
	@echo "Apply down migration 1"
	@$(MIGRATE) down 1

.PHONY: migrate-down-all
migrate-down-all:
	@echo "Apply all down migrations"
	@$(MIGRATE) down -all

.PHONY: migrate-force
migrate-force:
	@read -p "Enter force version: " version; \
	$(MIGRATE) force $${version}

.PHONY: migrate-version
migrate-version:
	@echo "Print current migration version"
	@$(MIGRATE) version

.PHONY: migrate-refresh
migrate-refresh:
	@echo "Refresh database migrations"
	@$(MIGRATE) down -all
	@$(MIGRATE) up