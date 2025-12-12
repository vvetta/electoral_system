include .env
export $(shell sed 's/=.*//' .env)

DB_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# DSN для тестовой базы (контейнер на localhost:PG_TEST_PORT)
TEST_DB_URL=postgres://$(PG_TEST_USER):$(PG_TEST_PASSWORD)@localhost:$(PG_TEST_PORT)/$(PG_TEST_DB_NAME)?sslmode=disable

MIGRATIONS_DIR=./migrations

# Создать новую пару файлов миграции: up/down
create-migration:
	@if [ -z "$(name)" ]; then echo "Usage: make create-migration name=add_users"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# Применить все миграции в основной БД
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Применить все миграции в ТЕСТОВОЙ БД (в контейнере)
migrate-test-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(TEST_DB_URL)" up

# Откатить одну миграцию вниз (основная БД)
migrate-down-1:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

# Откатить все (основная БД)
migrate-down-all:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

# Показать текущую версию (основная БД)
migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# Применить ровно N миграций вверх (основная БД)
migrate-up-n:
	@if [ -z "$(n)" ]; then echo "Usage: make migrate-up-n n=1"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up $(n)

# Форсировать версию (если миграции сломались и нужно выровнять метаданные, основная БД)
migrate-force:
	@if [ -z "$(v)" ]; then echo "Usage: make migrate-force v=3"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(v)


test-integration:
	@echo ">>> Запуск интеграционного тестирования..."
	@docker rm -f $(PG_TEST_CONTAINER_NAME) >/dev/null 2>&1 || true
	@docker run -d --name $(PG_TEST_CONTAINER_NAME) \
		-e POSTGRES_USER=$(PG_TEST_USER) \
		-e POSTGRES_PASSWORD=$(PG_TEST_PASSWORD) \
		-e POSTGRES_DB=$(PG_TEST_DB_NAME) \
		-p $(PG_TEST_PORT):5432 \
		postgres:latest

	@echo ">>> Ожидание запуска базы данных..."
	@until docker exec $(PG_TEST_CONTAINER_NAME) pg_isready -U $(PG_TEST_USER) >/dev/null 2>&1; do \
		sleep 1; \
	done

	@echo ">>> Применяем миграции к тестовой БД..."
	@if ! migrate -path $(MIGRATIONS_DIR) -database "$(TEST_DB_URL)" up; then \
		echo ">>> Миграции упали. Останавливаем контейнер..."; \
		docker stop $(PG_TEST_CONTAINER_NAME) >/dev/null 2>&1 || true; \
		docker rm $(PG_TEST_CONTAINER_NAME) >/dev/null 2>&1 || true; \
		exit 1; \
	fi

	@echo ">>> Запуск тестов..."
	@TEST_DATABASE_DSN="postgres://$(PG_TEST_USER):$(PG_TEST_PASSWORD)@localhost:$(PG_TEST_PORT)/$(PG_TEST_DB_NAME)?sslmode=disable" \
		go test ./... -v -integration ; \
		status=$$?; \
		echo ">>> Останавливаем PostgreSQL test container..."; \
		docker stop $(PG_TEST_CONTAINER_NAME) >/dev/null 2>&1 || true; \
		docker rm $(PG_TEST_CONTAINER_NAME) >/dev/null 2>&1 || true; \
		exit $$status
