# Terraform Provider for OpenShift Operators

[![GitHub release](https://img.shields.io/github/release/rh-mobb/terraform-provider-openshift.svg)](https://github.com/rh-mobb/terraform-provider-openshift/releases)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

A Terraform provider for managing OpenShift/OLM operators with simplified lifecycle management.

## Quick Start

```hcl
terraform {
  required_providers {
    openshift = {
      source  = "registry.terraform.io/rh-mobb/openshift"
      version = "~> 0.1"
    }
  }
}

provider "openshift" {
  kubeconfig = "~/.kube/config"
}

resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"
}
```

## Documentation

- [Provider Documentation](provider/README.md)
- [Resource Documentation](provider/docs/resources/openshift_operator.md)
- [Data Source Documentation](provider/docs/data-sources/openshift_operator.md)
- [Troubleshooting Guide](provider/docs/guides/troubleshooting.md)
- [Usage Examples](provider/USAGE.md)
- [Examples Directory](examples/)

## Best Practices

This provider follows Terraform community best practices. See [BEST_PRACTICES.md](BEST_PRACTICES.md) for details.

## Features

- **Single Resource**: Manage operator installation with a single `openshift_operator` resource
- **Automatic InstallPlan Approval**: Handles Manual approval automatically when version is pinned
- **Version Pinning**: Pin operators to specific versions with automatic Manual approval
- **CSV Waiting**: Automatically waits for ClusterServiceVersion to reach Succeeded phase
- **Lifecycle Management**: Handles Namespace, OperatorGroup, Subscription, InstallPlan, and CSV
- **Import Support**: Import existing operator installations into Terraform state
- **Update Support**: Handle version upgrades, channel changes, and configuration updates
- **Data Source**: Read operator status without managing it

## Requirements

- **Terraform**: >= 1.5
- **Go**: >= 1.21 (for building from source)
- **OpenShift/Kubernetes**: Access to a cluster with OLM installed

## Installation

### From Terraform Registry

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

### Local Development

See [provider/README.md](provider/README.md) for local development setup.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on contributing to this project.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/rh-mobb/terraform-provider-openshift/issues)
- **Security**: See [SECURITY.md](SECURITY.md) for security reporting

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes.
