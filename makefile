include ./.env.migration

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

MIGRATE_PATH=db/migrations
SEEDER_PATH=db/seeders/seeder.sql

print-test:
	@echo $(DB_URL) $(MIGRATE_PATH)

migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DB_URL) -path $(MIGRATE_PATH) up

migrate-down:
	migrate -database $(DB_URL) -path $(MIGRATE_PATH) down $(N)

db-seed:
	psql $(DB_URL) -f $(SEEDER_PATH)
	
# 	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f $(SEEDER_PATH)

db-reset:
	make migrate-down migrate-up db-seed


# 	migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) down -all
# 	migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) up
# 	PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB) -f $(SEEDER_PATH)