# Terraform Provider for OpenShift Operators

[![GitHub release](https://img.shields.io/github/release/rh-mobb/terraform-provider-openshift.svg)](https://github.com/rh-mobb/terraform-provider-openshift/releases)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rh-mobb/terraform-provider-openshift)](https://goreportcard.com/report/github.com/rh-mobb/terraform-provider-openshift)

A Terraform provider for managing OpenShift/OLM operators with simplified lifecycle management. This provider encapsulates the complexity of managing Namespace, OperatorGroup, Subscription, InstallPlan, and CSV resources into a single `openshift_operator` resource.

## Features

- **Single Resource**: Manage operator installation with a single `openshift_operator` resource
- **Automatic InstallPlan Approval**: Handles Manual approval automatically when version is pinned
- **Version Pinning**: Pin operators to specific versions with automatic Manual approval
- **CSV Waiting**: Automatically waits for ClusterServiceVersion to reach Succeeded phase
- **Lifecycle Management**: Handles Namespace, OperatorGroup, Subscription, InstallPlan, and CSV
- **Import Support**: Import existing operator installations into Terraform state
- **Update Support**: Handle version upgrades, channel changes, and configuration updates
- **Data Source**: Read operator status without managing it

## Installation

### From Terraform Registry

Add to your `required_providers` block:

```hcl
terraform {
  required_providers {
    openshift = {
      source  = "registry.terraform.io/rh-mobb/openshift"
      version = "~> 0.1"
    }
  }
}
```

Then run `terraform init` to download the provider.

### Local Development

```bash
cd provider
make install
```

This installs the provider to `~/.terraform.d/plugins/` for local use.

## Usage

### Basic Installation

```hcl
provider "openshift" {
  host     = var.api_url
  token    = var.k8s_token
  insecure = var.skip_tls_verify
}

resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"
}
```

### Version-Pinned Installation

```hcl
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"
  version   = "1.18.2"  # Automatically sets install_plan_approval to "Manual"
}
```

### With Custom Labels

```hcl
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"

  labels = {
    "app.kubernetes.io/managed-by" = "Terraform"
    "environment"                 = "production"
  }

  namespace_labels = {
    "openshift.io/cluster-monitoring" = "true"
  }
}
```

## Provider Configuration

The provider supports multiple authentication methods:

```hcl
provider "openshift" {
  # Option 1: Use kubeconfig file
  kubeconfig = "~/.kube/config"

  # Option 2: Use host and token
  host     = "https://api.example.com:6443"
  token    = var.k8s_token
  insecure = false

  # Option 3: Use default kubeconfig locations or in-cluster config
  # (no configuration needed)
}
```

### Provider Arguments

- `kubeconfig` (optional): Path to kubeconfig file. If not set, uses default locations or in-cluster config.
- `host` (optional): Kubernetes API server URL.
- `token` (optional, sensitive): Bearer token for API authentication.
- `insecure` (optional): Skip TLS certificate verification. Defaults to `false`.

## Resource Documentation

- [openshift_operator Resource](docs/resources/openshift_operator.md) - Manage operator installations
- [openshift_operator Data Source](docs/data-sources/openshift_operator.md) - Read operator status

## Examples

See the [examples](examples/) directory for more usage examples.

## Requirements

- **Terraform**: >= 1.5
- **Go**: >= 1.21 (for building from source)
- **OpenShift/Kubernetes**: Access to a cluster with OLM installed
- **Kubernetes API**: Version 1.20+ (for dynamic client support)

## Development

### Prerequisites

- Go 1.21 or later
- Terraform 1.5 or later
- Access to an OpenShift cluster (for integration testing)

### Building

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
# Unit tests
cd provider
make test

# Acceptance tests (requires cluster access)
TF_ACC=1 go test ./...
```

### Code Formatting

```bash
cd provider
make fmt
```

### Linting

```bash
cd provider
make lint
```

## Architecture

The provider encapsulates the complex operator installation flow:

1. **Namespace Creation** (if `create_namespace = true`)
2. **OperatorGroup Creation**
3. **Subscription Creation** with appropriate wait conditions
4. **InstallPlan Approval** (if Manual approval required)
5. **CSV Waiting** (if `wait_for_csv = true`)

All of this complexity is hidden from the user, who only needs to specify the operator name, channel, and source.

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines on contributing to this project.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](../LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/rh-mobb/terraform-provider-openshift/issues)
- **Security**: See [SECURITY.md](../SECURITY.md) for security reporting

## Roadmap

- [x] Basic CRUD operations
- [x] Import support
- [x] Update support for version/channel changes
- [ ] Operator uninstallation (CSV removal)
- [ ] Support for multiple catalog sources
- [ ] Enhanced error messages and diagnostics
- [ ] Comprehensive integration tests
