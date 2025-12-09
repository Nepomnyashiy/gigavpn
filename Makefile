# Makefile для автоматизации разработки GigaVPN

.PHONY: all build run stop clean test logs setup

all: build run

# === Go Backend ===
GO_BACKEND_DIR = ./backend-go
GO_APP_NAME = server

build-backend:
	@echo "--- Сборка Go бэкенда ---"
	@docker compose build gigavpn_backend

run-backend:
	@echo "--- Запуск Go бэкенда (без Docker Compose) ---"
	@cd $(GO_BACKEND_DIR) && go run ./cmd/$(GO_APP_NAME)

# === Python Bot ===
BOT_DIR = ./bot-python

build-bot:
	@echo "--- Сборка Python бота ---"
	@docker compose build bot

run-bot:
	@echo "--- Запуск Python бота (без Docker Compose, требуется BOT_TOKEN) ---"
	@cd $(BOT_DIR) && source .venv/bin/activate && python3 bot.py

# === Docker Compose ===
up:
	@echo "--- Запуск всех сервисов через Docker Compose ---"
	@docker compose up -d --build --force-recreate

down:
	@echo "--- Остановка всех сервисов Docker Compose ---"
	@docker compose down

logs:
	@echo "--- Просмотр логов всех сервисов Docker Compose ---"
	@docker compose logs -f

# === Development Utilities ===
test:
	@echo "--- Запуск тестов (пока заглушка) ---"
	@echo "Тесты еще не реализованы."

clean:
	@echo "--- Очистка (удаление логов, кэша) ---"
	@rm -f $(GO_BACKEND_DIR)/*.log
	@rm -f $(BOT_DIR)/*.log
	@rm -f $(GO_BACKEND_DIR)/$(GO_APP_NAME)
	@docker system prune -f --volumes

.PHONY: docker-logs
docker-logs:
	@docker compose logs -f

.PHONY: db-shell
db-shell:
	@docker exec -it gigavpn_postgres psql -U gigavpn_user -d gigavpn_db

.PHONY: db-migrations
db-migrations:
	@echo "--- Применение миграций к БД ---"
	@docker run --rm -v $(GO_BACKEND_DIR)/migrations:/migrations --network gigavpn_gigavpn_net migrate/migrate \
		-path=/migrations/ -database 'postgres://gigavpn_user:gigavpn_password_DoNotUseInProd@postgres:5432/gigavpn_db?sslmode=disable' up

.PHONY: db-migrate-down
db-migrate-down:
	@echo "--- Откат миграций БД ---"
	@docker run --rm -v $(GO_BACKEND_DIR)/migrations:/migrations --network gigavpn_gigavpn_net migrate/migrate \
		-path=/migrations/ -database 'postgres://gigavpn_user:gigavpn_password_DoNotUseInProd@postgres:5432/gigavpn_db?sslmode=disable' down 1

.PHONY: db-migrate-status
db-migrate-status:
	@echo "--- Статус миграций БД ---"
	@docker run --rm -v $(GO_BACKEND_DIR)/migrations:/migrations --network gigavpn_gigavpn_net migrate/migrate \
		-path=/migrations/ -database 'postgres://gigavpn_user:gigavpn_password_DoNotUseInProd@postgres:5432/gigavpn_db?sslmode=disable' status

# === Ansible ===
ANSIBLE_PLAYBOOKS_DIR = ./infrastructure/ansible/playbooks
ANSIBLE_INVENTORY_DIR = ./infrastructure/ansible/inventory

ansible-deploy:
	@echo "--- Запуск Ansible плейбука для развертывания VPN-узлов ---"
	ansible-playbook -i $(ANSIBLE_INVENTORY_DIR)/hosts.ini $(ANSIBLE_PLAYBOOKS_DIR)/site.yml
