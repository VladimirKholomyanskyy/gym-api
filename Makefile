# Simple Makefile for a Go project
OPENAPI_REPO=git@github.com:VladimirKholomyanskyy/gym-api-contracts.git
OPENAPI_DIR=api-contracts
GENERATOR_IMAGE=openapitools/openapi-generator-cli
OPENAPI_SPEC=$(OPENAPI_DIR)/openapi/v1/spec.yaml
OPENAPI_OUTPUT=internal/api

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@docker compose up --build

# Shutdown DB container
docker-down:
	@docker compose down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integration Tests
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main.exe

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/cosmtrek/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

# Clone or update OpenAPI specification
openapi-update:
	@if [ ! -d "$(OPENAPI_DIR)" ]; then \
		echo "Cloning OpenAPI repo..."; \
		git clone $(OPENAPI_REPO) $(OPENAPI_DIR); \
	else \
		echo "Updating OpenAPI repo..."; \
		git -C $(OPENAPI_DIR) pull; \
	fi

# Run OpenAPI code generation
openapi-generate:
	rm -rf $(OPENAPI_OUTPUT)  # Remove old generated code
	docker run --rm -v $(PWD):/local $(GENERATOR_IMAGE) generate \
		-i /local/$(OPENAPI_SPEC) -g go-server -o /local/$(OPENAPI_OUTPUT) \
		--git-repo-id gym-api --git-user-id VladimirKholomyanskyy \
		--additional-properties=onlyInterfaces=true,outputAsLibrary=true

# Run OpenAPI update & generate
openapi-all: openapi-update openapi-generate

.PHONY: all build run test clean watch docker-run docker-down itest openapi-update openapi-generate openapi-all
