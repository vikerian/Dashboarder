.PHONY: all build clean run-local docker-build docker-up docker-down test

# Variables
DOCKER_COMPOSE = docker-compose
GO = go

# Default target
all: build

# Build all services locally
build:
	@echo "Building all services..."
	$(GO) build -o bin/scraper-mmdecin ./services/scraper-mmdecin
	$(GO) build -o bin/scraper-traffic ./services/scraper-traffic
	$(GO) build -o bin/data-router ./services/data-router
	$(GO) build -o bin/api-server ./services/api-server
	$(GO) build -o bin/web-server ./services/web-server

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/

# Run services locally (requires local MongoDB, MQTT, SiriDB)
run-local: build
	@echo "Starting services locally..."
	@echo "Make sure MongoDB, Mosquitto, and SiriDB are running!"
	./bin/api-server &
	./bin/web-server &
	@echo "Services started. API on :8080, Web on :8090"

# Docker commands
docker-build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

docker-up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d

docker-down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	$(GO) mod download
	$(GO) mod tidy
	mkdir -p mosquitto/config mosquitto/data mosquitto/log
	echo "listener 1883" > mosquitto/config/mosquitto.conf
	echo "allow_anonymous true" >> mosquitto/config/mosquitto.conf

# Test
test:
	@echo "Running tests..."
	$(GO) test ./...

# Run specific service
run-mmdecin:
	$(GO) run ./services/scraper-mmdecin/main.go

run-traffic:
	$(GO) run ./services/scraper-traffic/main.go

run-router:
	$(GO) run ./services/data-router/main.go

run-api:
	$(GO) run ./services/api-server/main.go

run-web:
	$(GO) run ./services/web-server/main.go

# IoT simulator for testing
run-iot-simulator:
	$(GO) run ./tools/iot-simulator/main.go

# Help
help:
	@echo "Available targets:"
	@echo "  make build          - Build all services"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-build   - Build Docker images"
	@echo "  make docker-up      - Start all services in Docker"
	@echo "  make docker-down    - Stop all Docker services"
	@echo "  make docker-logs    - View Docker logs"
	@echo "  make dev-setup      - Setup development environment"
	@echo "  make test           - Run tests"
	@echo "  make run-local      - Run services locally"
	@echo "  make run-<service>  - Run specific service"
