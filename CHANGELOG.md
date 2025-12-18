# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.1] - 2025-12-18

### Fixed
- Fixed golangci-lint configuration compatibility issues
- Fixed all linting errors (errcheck, revive, staticcheck, unused)
- Fixed TLSClientConfig field access warnings
- Added missing package comments for all packages
- Fixed unused variables and parameters
- Removed duplicate version.go file

### Changed
- Removed GPG signing requirement from release workflow (checksums only)
- Updated release documentation to remove GPG prerequisites
- Improved code quality and maintainability

## [0.1.0] - 2025-12-18

### Added
- Initial release of the OpenShift Operator Terraform Provider
- Support for installing OpenShift operators via OLM
- Automatic namespace creation
- OperatorGroup management
- Subscription creation with channel and source configuration
- Automatic InstallPlan approval for Manual approval strategy
- CSV waiting and phase tracking
- Version pinning support
- Comprehensive resource attributes for operator lifecycle management
- Basic CRUD operations for `openshift_operator` resource
- Support for Automatic and Manual InstallPlan approval
- Namespace and OperatorGroup creation
- CSV phase tracking
- Resource import functionality
- Update support for version/channel changes
- Data source for reading operator status
- Comprehensive documentation
- Examples directory with multiple use cases
- Website documentation for Terraform Registry
- Pre-commit hooks and code quality tools
- CI/CD workflows for testing and releases
- Dependabot for dependency management

### Changed
- Provider namespace set to `rh-mobb/openshift` for Terraform Registry publication

[Unreleased]: https://github.com/rh-mobb/terraform-provider-openshift/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/rh-mobb/terraform-provider-openshift/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/rh-mobb/terraform-provider-openshift/releases/tag/v0.1.0
