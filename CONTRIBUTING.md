# Contributing to Terraform Provider for OpenShift Operators

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to this project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/terraform-provider-openshift.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit your changes: `git commit -m "Add some feature"`
7. Push to your fork: `git push origin feature/your-feature-name`
8. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.21 or later
- Terraform 1.5 or later
- Access to an OpenShift cluster (for integration testing)

### Building the Provider

```bash
cd provider
make build
```

### Installing Locally

```bash
cd provider
make install
```

### Running Tests

```bash
cd provider
make test
```

## Code Style

- Follow Go best practices and conventions
- Use `gofmt` to format code: `make fmt`
- Run linter: `make lint`
- Add comments for exported functions and types
- Keep functions focused and small

## Testing Guidelines

### Unit Tests

- Write unit tests for all new functionality
- Mock external dependencies (Kubernetes API)
- Aim for high test coverage

### Integration Tests

- Integration tests require a real OpenShift cluster
- Use acceptance tests for full resource lifecycle testing
- Clean up test resources after tests complete

### Running Tests

```bash
# Unit tests only
go test ./...

# With coverage
go test -coverprofile=coverage.out ./...

# Acceptance tests (requires cluster access)
TF_ACC=1 go test ./...
```

## Commit Messages

- Use clear, descriptive commit messages
- Reference issue numbers when applicable: `Fixes #123`
- Follow conventional commit format when possible:
  - `feat: Add new feature`
  - `fix: Fix bug in update logic`
  - `docs: Update README`
  - `test: Add unit tests`

## Pull Request Process

1. Ensure all tests pass
2. Update documentation if needed
3. Update CHANGELOG.md if applicable
4. Ensure your PR description clearly describes the changes
5. Request review from maintainers

## Reporting Issues

When reporting issues, please include:

- Provider version
- Terraform version
- OpenShift/Kubernetes version
- Terraform configuration (sanitized)
- Error messages and logs
- Steps to reproduce

## Code Review

- All code must be reviewed before merging
- Address review comments promptly
- Be respectful and constructive in reviews

## Questions?

Feel free to open an issue for questions or reach out to the maintainers.

Thank you for contributing!
