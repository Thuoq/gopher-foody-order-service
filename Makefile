# ====== CONFIG ======
APP_NAME=gopher-foody-order
MIGRATIONS_PATH=./migrations

# dùng host.docker.internal cho Mac để kết nối tới port đã export ở docker-compose
DB_URL?=postgres://postgres:thuong123@host.docker.internal:5435/foody_order_dev_db?sslmode=disable

# image migrate
MIGRATE_IMAGE=migrate/migrate

# ====== MIGRATION ======

migrate-up:
	docker run --rm \
	-v $(PWD)/migrations:/migrations \
	$(MIGRATE_IMAGE) \
	-path=/migrations \
	-database "$(DB_URL)" up

migrate-down:
	docker run --rm \
	-v $(PWD)/migrations:/migrations \
	$(MIGRATE_IMAGE) \
	-path=/migrations \
	-database "$(DB_URL)" down 1

migrate-down-all:
	docker run --rm \
	-v $(PWD)/migrations:/migrations \
	$(MIGRATE_IMAGE) \
	-path=/migrations \
	-database "$(DB_URL)" down -all

migrate-reset: migrate-down-all migrate-up

migrate-force:
	docker run --rm \
	-v $(PWD)/migrations:/migrations \
	$(MIGRATE_IMAGE) \
	-path=/migrations \
	-database "$(DB_URL)" force $(v)

migrate-version:
	docker run --rm \
	-v $(PWD)/migrations:/migrations \
	$(MIGRATE_IMAGE) \
	-path=/migrations \
	-database "$(DB_URL)" version

migrate-create:
	docker run --rm \
	-v $(PWD)/migrations:/migrations \
	$(MIGRATE_IMAGE) \
	create -ext sql -dir /migrations -seq $(name)

# ====== APP ======

run:
	go run cmd/server/main.go

build:
	go build -o bin/$(APP_NAME) cmd/server/main.go

# ====== CLEAN ======

clean:
	rm -rf bin/