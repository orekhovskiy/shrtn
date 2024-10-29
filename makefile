# Variables
APP_NAME = shrtn
DOCKER_COMPOSE_FILE = ./development/docker-compose.yaml
DATABASE_DSN = postgres://postgres:password@localhost:5432/shrtn
GO_FILE = ./cmd/shortener/main.go

# Targets

# Start Docker Compose for the PostgreSQL service
.PHONY: up
up:
	@echo "Starting Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Stop Docker Compose
.PHONY: down
down:
	@echo "Stopping Docker Compose..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Build the Go application
.PHONY: build
build:
	@echo "Building the application..."
	go build -o $(APP_NAME) $(GO_FILE)

# Run the Go application
.PHONY: run
run: build
	@echo "Running the application..."
	DATABASE_DSN=$(DATABASE_DSN) ./$(APP_NAME)

# Clean up generated files
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)

# Format Go files
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Lint Go code
.PHONY: lint
lint:
	@echo "Linting Go code..."
	golangci-lint run

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Full workflow: start docker, build, and run
.PHONY: all
all: up run
