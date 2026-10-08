include ./.env.migration
include ./.env.redis


REDIS_URL = redis://$(REDIS_USER):$(REDIS_PASS)@$(REDIS_HOST):$(REDIS_PORT)
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
	
db-reset:
	make migrate-down migrate-up db-seed flushdb

redis-cli:
	@redis-cli -h $(REDIS_HOST) -p $(REDIS_PORT) --user $(REDIS_USER) --pass "$(REDIS_PASS)"

flushdb:
	@redis-cli -h $(REDIS_HOST) -p $(REDIS_PORT) --user $(REDIS_USER) --pass "$(REDIS_PASS)" FLUSHDB
	@echo "Redis aktif berhasil dibersihkan!"

flushall:
	@redis-cli -h $(REDIS_HOST) -p $(REDIS_PORT) --user $(REDIS_USER) --pass "$(REDIS_PASS)" FLUSHALL
	@echo "Semua database Redis berhasil dibersihkan!"

keys:
	@redis-cli -h $(REDIS_HOST) -p $(REDIS_PORT) --user $(REDIS_USER) --pass "$(REDIS_PASS)" KEYS "*"