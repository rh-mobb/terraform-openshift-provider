# Examples

This directory contains example Terraform configurations demonstrating various use cases for the OpenShift provider.

## Examples

### Basic

[`basic/`](./basic/) - Simple operator installation example.

```bash
cd examples/basic
terraform init
terraform plan
terraform apply
```

### Version Pinned

[`version-pinned/`](./version-pinned/) - Install operator with version pinning.

```bash
cd examples/version-pinned
terraform init
terraform plan
terraform apply
```

### Multiple Operators

[`multiple-operators/`](./multiple-operators/) - Install multiple operators.

```bash
cd examples/multiple-operators
terraform init
terraform plan
terraform apply
```

### Data Source

[`data-source/`](./data-source/) - Use data source to read operator status.

```bash
cd examples/data-source
terraform init
terraform plan
terraform apply
```

## Prerequisites

- Terraform >= 1.5
- Access to an OpenShift cluster
- Valid kubeconfig or cluster credentials
- OLM (Operator Lifecycle Manager) installed on the cluster

## Notes

- Examples use placeholder values - update with your actual cluster details
- Some examples may require specific operator channels or versions
- Always review and test in a non-production environment first
