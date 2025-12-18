.PHONY: build install test clean version

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS = -X github.com/redhat/terraform-provider-openshift-operator/internal/version.Version=$(VERSION) \
          -X github.com/redhat/terraform-provider-openshift-operator/internal/version.GitCommit=$(GIT_COMMIT) \
          -X github.com/redhat/terraform-provider-openshift-operator/internal/version.BuildDate=$(BUILD_DATE)

# Build the provider
build:
	go build -ldflags "$(LDFLAGS)" -o terraform-provider-openshift

# Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Install the provider to local Terraform plugins directory
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/rh-mobb/openshift/0.1.0/$$(go env GOOS)_$$(go env GOARCH)
	cp terraform-provider-openshift ~/.terraform.d/plugins/registry.terraform.io/rh-mobb/openshift/0.1.0/$$(go env GOOS)_$$(go env GOARCH)/

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f terraform-provider-openshift

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy
