# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BINARY_NAME=nous_backend
BINARY_UNIX=$(BINARY_NAME)_unix

# Air
AIR_CMD=air

# Templ
TEMPL_CMD=templ

# Main package path
MAIN_PATH=cmd/api/main.go

# Build the project
all: generate build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Generate templ files
generate:
	$(TEMPL_CMD) generate

# Run the project
run: generate
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)
	./$(BINARY_NAME)

# Clean the binary
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test: generate
	$(GOTEST) -v ./...

dev: generate
	$(AIR_CMD) -c .air.toml

# Cross compilation for Linux
build-linux: generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

.PHONY: all build generate run clean test dev build-linux
