# Build the application
all: build

build:
	@echo "Building..."
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Clear docker volume
db.clear:
	@echo "Clearing DB volume..."
	@docker volume rm transaction-routine_psql_volume

# Create DB container
db.up:
	@echo "Starting DB container..."
	@docker compose up psql -d

# Shutdown DB container
db.down:
	@echo "Stopping DB container..."
	@docker compose down psql

# Run migrations
migrate.up:
	@echo "Running migrations..."
	@docker compose up migrate -d

# Reset the database
db.reset:
	@echo "Resetting the database..."
	@make db.down
	@make db.clear
	@make db.up
	@sleep 1
	@make migrate.up

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Generate mocks
mocks:
	@echo "Generating mocks..."
	@go generate ./...

# Run integration/load tests
loadtest:
	@docker run --rm -i --net=host grafana/k6 run - <k6/script.js

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean