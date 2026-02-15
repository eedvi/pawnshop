# Makefile for Pawnshop API

.PHONY: all build run test clean migrate docker help

# Variables
APP_NAME=pawnshop
BUILD_DIR=./build
MAIN_PATH=./cmd/api/main.go
DOCKER_IMAGE=pawnshop-api

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w"

all: build

## build: Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

## run: Run the application
run:
	@echo "Running $(APP_NAME)..."
	$(GORUN) $(MAIN_PATH)

## dev: Run with hot reload (requires air)
dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -cover ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: Run linter
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

## tidy: Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

## download: Download dependencies
download:
	@echo "Downloading dependencies..."
	$(GOMOD) download

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## migrate-up: Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	@which migrate > /dev/null || (echo "Installing migrate..." && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pawnshop?sslmode=disable" up

## migrate-down: Run database migrations down
migrate-down:
	@echo "Running migrations down..."
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pawnshop?sslmode=disable" down

## migrate-create: Create new migration (usage: make migrate-create name=create_table)
migrate-create:
	@echo "Creating migration $(name)..."
	migrate create -ext sql -dir ./migrations -seq $(name)

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) -f docker/Dockerfile .

## docker-run: Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

## docker-compose-up: Start all services with docker-compose
docker-compose-up:
	@echo "Starting services..."
	docker-compose -f docker/docker-compose.yml up -d

## docker-compose-down: Stop all services
docker-compose-down:
	@echo "Stopping services..."
	docker-compose -f docker/docker-compose.yml down

## swagger: Generate Swagger docs
swagger:
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/api/main.go -o api/docs

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
