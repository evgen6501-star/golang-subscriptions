include .env
export

PROJECT_ROOT := $(shell pwd)


run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	go mod tidy && \
	go run cmd/main.go


build:
	docker build --network=host -t golang-subscriptions-sub-app .

up: build
	docker compose up -d
	@echo "✅ Сервисы запущены"
	@echo "📌 Приложение: http://localhost:8080"
	@echo "📌 Логи: docker compose logs -f sub-app"

down:
	docker compose down

logs:
	docker compose logs -f

clean:
	docker compose down -v
	sudo rm -rf out/pgdata