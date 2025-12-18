.PHONY: help build test install clean fmt lint version deps verify acceptance-test

# Default target
.DEFAULT_GOAL := help

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS = -X github.com/redhat/terraform-provider-openshift-operator/internal/version.Version=$(VERSION) \
          -X github.com/redhat/terraform-provider-openshift-operator/internal/version.GitCommit=$(GIT_COMMIT) \
          -X github.com/redhat/terraform-provider-openshift-operator/internal/version.BuildDate=$(BUILD_DATE)

# Provider directory
PROVIDER_DIR = provider
BINARY_NAME = terraform-provider-openshift

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the provider binary
	@echo "Building $(BINARY_NAME)..."
	cd $(PROVIDER_DIR) && go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)
	@echo "Build complete: $(PROVIDER_DIR)/$(BINARY_NAME)"

test: ## Run unit tests
	@echo "Running unit tests..."
	cd $(PROVIDER_DIR) && go test -v -coverprofile=coverage.out ./...
	@echo "Tests complete"

test-coverage: test ## Run tests and show coverage
	cd $(PROVIDER_DIR) && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: $(PROVIDER_DIR)/coverage.html"

acceptance-test: ## Run acceptance tests (requires TF_ACC=1 and cluster access)
	@if [ -z "$$TF_ACC" ]; then \
		echo "Error: TF_ACC=1 must be set for acceptance tests"; \
		echo "Usage: TF_ACC=1 make acceptance-test"; \
		exit 1; \
	fi
	@echo "Running acceptance tests..."
	cd $(PROVIDER_DIR) && go test -v -tags=acceptance ./...

install: build ## Install provider to local Terraform plugins directory
	@echo "Installing provider..."
	@mkdir -p ~/.terraform.d/plugins/registry.terraform.io/rh-mobb/openshift/$(VERSION)/$$(go env GOOS)_$$(go env GOARCH)
	@cp $(PROVIDER_DIR)/$(BINARY_NAME) ~/.terraform.d/plugins/registry.terraform.io/rh-mobb/openshift/$(VERSION)/$$(go env GOOS)_$$(go env GOARCH)/
	@echo "Provider installed to ~/.terraform.d/plugins/registry.terraform.io/rh-mobb/openshift/$(VERSION)/$$(go env GOOS)_$$(go env GOARCH)/"

fmt: ## Format Go code
	@echo "Formatting code..."
	cd $(PROVIDER_DIR) && go fmt ./...
	@echo "Formatting complete"

lint: ## Run golangci-lint
	@echo "Running linter..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint not found. Installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest; \
	fi
	cd $(PROVIDER_DIR) && golangci-lint run
	@echo "Linting complete"

lint-install: ## Install golangci-lint
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest
	@echo "golangci-lint installed"

deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	cd $(PROVIDER_DIR) && go mod download
	cd $(PROVIDER_DIR) && go mod tidy
	@echo "Dependencies updated"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(PROVIDER_DIR)/$(BINARY_NAME)
	rm -f $(PROVIDER_DIR)/$(BINARY_NAME)-*
	rm -f $(PROVIDER_DIR)/coverage.out
	rm -f $(PROVIDER_DIR)/coverage.html
	rm -rf $(PROVIDER_DIR)/dist
	@echo "Clean complete"

version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@if [ -f $(PROVIDER_DIR)/$(BINARY_NAME) ]; then \
		$(PROVIDER_DIR)/$(BINARY_NAME) -version 2>&1 || echo "Version command not available"; \
	fi

verify: fmt lint test ## Run all verification checks (format, lint, test)
	@echo "All verification checks passed"

# Build for multiple platforms (for releases)
build-all: ## Build binaries for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(PROVIDER_DIR)/dist
	cd $(PROVIDER_DIR) && \
		GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_$(VERSION)_linux_amd64 && \
		GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_$(VERSION)_linux_arm64 && \
		GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_$(VERSION)_darwin_amd64 && \
		GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_$(VERSION)_darwin_arm64 && \
		GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_$(VERSION)_windows_amd64.exe
	@echo "Build complete. Binaries in $(PROVIDER_DIR)/dist/"

checksums: build-all ## Generate SHA256 checksums for binaries
	@echo "Generating checksums..."
	cd $(PROVIDER_DIR)/dist && sha256sum $(BINARY_NAME)_* > $(BINARY_NAME)_$(VERSION)_SHA256SUMS
	@echo "Checksums generated: $(PROVIDER_DIR)/dist/$(BINARY_NAME)_$(VERSION)_SHA256SUMS"

# CI/CD helpers
ci-test: deps test ## Run tests for CI (downloads deps first)
	@echo "CI tests complete"

ci-lint: deps lint ## Run linting for CI (downloads deps first)
	@echo "CI linting complete"

ci-build: deps build ## Build for CI (downloads deps first)
	@echo "CI build complete"
