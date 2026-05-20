SERVICE_NAME := customer-service
SERVICE_DIR := services/customer-service
SERVICE_CMD := ./$(SERVICE_DIR)/cmd/customer-service
BINARY_DIR := bin
BINARY_PATH := $(BINARY_DIR)/$(SERVICE_NAME)

CONFIG_PATH ?= configs/customer.env
DOCKER_ENV_FILE ?= configs/customer.docker.env
IMAGE ?= softplace/customer-service:latest
COMPOSE_FILE ?= docker-compose.customer.yml

.PHONY: help customer-run customer-build customer-test customer-test-race customer-clean \
	customer-docker-build customer-docker-run customer-docker-stop \
	customer-compose-up customer-compose-down customer-compose-logs customer-compose-migrate

help:
	@echo "Available targets:"
	@echo "  customer-run           Run customer-service locally"
	@echo "  customer-build         Build customer-service binary"
	@echo "  customer-test          Run customer-service tests"
	@echo "  customer-test-race     Run customer-service tests with race detector"
	@echo "  customer-clean         Remove built binaries"
	@echo "  customer-docker-build  Build customer-service Docker image"
	@echo "  customer-docker-run    Run customer-service container"
	@echo "  customer-docker-stop   Stop customer-service container"
	@echo "  customer-compose-up    Start customer-service + postgres via docker compose"
	@echo "  customer-compose-migrate  Apply SQL migrations to postgres in compose"
	@echo "  customer-compose-down  Stop docker compose services"
	@echo "  customer-compose-logs  Tail docker compose logs"

customer-run:
	CONFIG_PATH=$(CONFIG_PATH) go run $(SERVICE_CMD)

customer-build:
	mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_PATH) $(SERVICE_CMD)

customer-test:
	go test -count=1 ./$(SERVICE_DIR)/...

customer-test-race:
	go test -race -count=1 ./$(SERVICE_DIR)/...

customer-clean:
	rm -rf $(BINARY_DIR)

customer-docker-build:
	docker build -f $(SERVICE_DIR)/Dockerfile -t $(IMAGE) .

customer-docker-run:
	docker run --rm -d \
		--name $(SERVICE_NAME) \
		--env-file $(DOCKER_ENV_FILE) \
		-e CONFIG_PATH=/app/configs/customer.docker.env \
		-p 8442:8442 \
		$(IMAGE)

customer-docker-stop:
	-docker stop $(SERVICE_NAME)

customer-compose-up:
	docker compose -f $(COMPOSE_FILE) up -d --build

customer-compose-migrate:
	docker compose -f $(COMPOSE_FILE) up -d softplace-db
	@echo "Applying migrations from ./migrations to softplace-db..."
	@set -e; \
	for file in $$(ls -1 migrations/*.sql | sort); do \
		echo "-> $$file"; \
		docker compose -f $(COMPOSE_FILE) exec -T softplace-db \
			psql -U postgres -d softplace -v ON_ERROR_STOP=1 -f /migrations/$$(basename $$file); \
	done

customer-compose-down:
	docker compose -f $(COMPOSE_FILE) down -v

customer-compose-logs:
	docker compose -f $(COMPOSE_FILE) logs -f

