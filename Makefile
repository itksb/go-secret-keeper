# Makefile for the project
# Make sure that the stating lines are indented with a tab character, not spaces.

SHELL := /bin/bash
PWD := $(shell pwd)
GIT_BUILD_COMMIT := $(shell git rev-parse HEAD)
GIT_BUILD_TAG := $(shell git tag -l --contains HEAD | head -1)
OUTPUT_DIR := ./output
BUILD_DATE_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
CLIENT_BINARY_NAME := $(OUTPUT_DIR)/client
SERVER_BINARY_NAME := $(OUTPUT_DIR)/server

all: help

.PHONY: build-client
build-client: ## Build client binary
	@echo "Building client binary"
	@mkdir -p $(OUTPUT_DIR)
	@go build -ldflags "-X main.gitBuildCommit=$(GIT_BUILD_COMMIT) -X main.gitBuildTag=$(GIT_BUILD_TAG) -X main.buildDateTime=$(BUILD_DATE_TIME) " -o $(CLIENT_BINARY_NAME) ./cmd/client
	@echo "output: $(PWD)/$(CLIENT_BINARY_NAME)"

.PHONY: build-server
build-server: ## Build server binary
	@echo "Building server binary"
	@mkdir -p $(OUTPUT_DIR)
	@go build -ldflags "-X main.gitBuildCommit=$(GIT_BUILD_COMMIT) -X main.gitBuildTag=$(GIT_BUILD_TAG) -X main.buildDateTime=$(BUILD_DATE_TIME) " -o $(SERVER_BINARY_NAME) ./cmd/server
	@echo "output: $(PWD)/$(SERVER_BINARY_NAME)"

.PHONY: build-all
build-all: ## Build all binaries: server and client
	@echo "Building all binaries"
	@mkdir -p $(OUTPUT_DIR)
	@$(MAKE) build-client
	@$(MAKE) build-server


.PHONY: fmt
fmt: ## Format the source code
	@echo "Formatting the source code"
	go fmt ./...

.PHONY: lint
vet: ## Vet the source code
	@echo "Linting the source code"
	go vet ./...

.PHONY: test
test: ## Run tests
	@echo "Running tests"
	go test -v ./...

.PHONY: cover-client
cover-client: ## Run tests with coverage
	@echo "Running tests with coverage"
	go test -coverprofile=clientprofile.out ./internal/client/...
	go tool cover -html=clientprofile.out

.PHONY: cover-server
cover-server: ## Run tests with coverage
	@echo "Running tests with coverage"
	go test -coverprofile=serverprofile.out ./internal/server/...
	go tool cover -html=serverprofile.out

.PHONY: help
help: ## Show current help message
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' ./Makefile | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}'
