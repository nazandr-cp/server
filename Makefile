# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server
BINARY_UNIX=$(BINARY_NAME)_unix

# Main package path
MAIN_PATH=./cmd/server

.PHONY: all build clean test coverage deps help

all: test build

## build: Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

## clean: Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

## test: Run tests
test:
	$(GOTEST) -v ./...

## coverage: Run tests with coverage
coverage:
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out

## deps: Download and verify dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

## tidy: Clean up dependencies
tidy:
	$(GOMOD) tidy

## run: Run the application
run:
	$(GOCMD) run $(MAIN_PATH)

## build-linux: Cross compilation for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

## docker-build: Build Docker image
docker-build:
	docker build -t go-server .

## docker-run: Run Docker container
docker-run:
	docker-compose up --build

## lint: Run linter
lint:
	golangci-lint run

## help: Show this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'