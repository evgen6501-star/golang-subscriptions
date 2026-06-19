include .env
export


export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up -d sub-postgres

env-down:
	docker compose down sub-postgres

env-cleanup:
	@read -p "Очистить все файлы???? y/N: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down -v sub-postgres && \
		sudo rm -rf  out/pgdata && \
		echo "удалено"; \
	else \
		echo "Отмена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
	echo "отсутствует параметр seq"; \
	exit 1; \
	fi; \

	docker compose run --rm sub-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"


migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-down:
	@if [ -z "$(action)" ]; then \
	echo "отсутствует параметр action"; \
	exit 1; \
	fi; \
	docker compose run --rm sub-postgres-migrate \
	-path /migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@sub-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	"$(action)"