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