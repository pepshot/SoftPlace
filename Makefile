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
	customer-compose-up customer-compose-down customer-compose-logs customer-compose-migrate \
	supplier-run supplier-build supplier-test supplier-test-race supplier-clean \
	supplier-docker-build supplier-docker-run supplier-docker-stop \
	supplier-compose-up supplier-compose-down supplier-compose-logs supplier-compose-migrate

help:
	@echo "Available targets:"
	@echo ""
	@echo "Customer Service:"
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
	@echo ""
	@echo "Supplier Service:"
	@echo "  supplier-run           Run supplier-service locally"
	@echo "  supplier-build         Build supplier-service binary"
	@echo "  supplier-test          Run supplier-service tests"
	@echo "  supplier-test-race     Run supplier-service tests with race detector"
	@echo "  supplier-clean         Remove built binaries"
	@echo "  supplier-docker-build  Build supplier-service Docker image"
	@echo "  supplier-docker-run    Run supplier-service container"
	@echo "  supplier-docker-stop   Stop supplier-service container"
	@echo "  supplier-compose-up    Start supplier-service + postgres via docker compose"
	@echo "  supplier-compose-migrate  Apply SQL migrations to postgres in compose"
	@echo "  supplier-compose-down  Stop docker compose services"
	@echo "  supplier-compose-logs  Tail docker compose logs"

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

# Supplier Service Targets
SUPPLIER_SERVICE_NAME := supplier-service
SUPPLIER_SERVICE_DIR := services/supplier-service
SUPPLIER_SERVICE_CMD := ./$(SUPPLIER_SERVICE_DIR)/cmd/supplier-service
SUPPLIER_BINARY_PATH := $(BINARY_DIR)/$(SUPPLIER_SERVICE_NAME)

SUPPLIER_CONFIG_PATH ?= configs/supplier.env
SUPPLIER_DOCKER_ENV_FILE ?= configs/supplier.docker.env
SUPPLIER_IMAGE ?= softplace/supplier-service:latest
SUPPLIER_COMPOSE_FILE ?= docker-compose.supplier.yml

supplier-run:
	CONFIG_PATH=$(SUPPLIER_CONFIG_PATH) go run $(SUPPLIER_SERVICE_CMD)

supplier-build:
	mkdir -p $(BINARY_DIR)
	go build -o $(SUPPLIER_BINARY_PATH) $(SUPPLIER_SERVICE_CMD)

supplier-test:
	go test -count=1 ./$(SUPPLIER_SERVICE_DIR)/...

supplier-test-race:
	go test -race -count=1 ./$(SUPPLIER_SERVICE_DIR)/...

supplier-clean:
	rm -rf $(BINARY_DIR)

supplier-docker-build:
	docker build -f $(SUPPLIER_SERVICE_DIR)/Dockerfile -t $(SUPPLIER_IMAGE) .

supplier-docker-run:
	docker run --rm -d \
		--name $(SUPPLIER_SERVICE_NAME) \
		--env-file $(SUPPLIER_DOCKER_ENV_FILE) \
		-e CONFIG_PATH=/app/configs/supplier.docker.env \
		-p 8443:8443 \
		-p 9091:9091 \
		$(SUPPLIER_IMAGE)

supplier-docker-stop:
	-docker stop $(SUPPLIER_SERVICE_NAME)

supplier-compose-up:
	docker compose -f $(SUPPLIER_COMPOSE_FILE) up -d --build

supplier-compose-migrate:
	docker compose -f $(SUPPLIER_COMPOSE_FILE) up -d softplace-db
	@echo "Applying migrations from ./migrations to softplace-db..."
	@set -e; \
	for file in $$(ls -1 migrations/*.sql | sort); do \
		echo "-> $$file"; \
		docker compose -f $(SUPPLIER_COMPOSE_FILE) exec -T softplace-db \
			psql -U postgres -d softplace -v ON_ERROR_STOP=1 -f /migrations/$$(basename $$file); \
	done

supplier-compose-down:
	docker compose -f $(SUPPLIER_COMPOSE_FILE) down -v

supplier-compose-logs:
	docker compose -f $(SUPPLIER_COMPOSE_FILE) logs -f

