# Using the OpenShift Operator Provider

## Installation

### Local Development

1. Build the provider:
```bash
cd provider
make build
```

2. Install to local Terraform plugins directory:
```bash
make install
```

This installs the provider to `~/.terraform.d/plugins/registry.terraform.io/rh-mobb/openshift/0.1.0/$(go env GOOS)_$(go env GOARCH)/`

## Configuration

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    openshift = {
      source  = "registry.terraform.io/rh-mobb/openshift"
      version = "0.1.0"
    }
  }
}

provider "openshift" {
  host     = var.api_url
  token    = var.k8s_token
  insecure = var.skip_tls_verify
}
```

## Example: Installing GitOps Operator

```hcl
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"
  version   = "1.18.2"  # Pins version and sets Manual approval automatically

  labels = {
    "app.kubernetes.io/managed-by" = "Terraform"
  }

  namespace_labels = {
    "openshift.io/cluster-monitoring" = "true"
  }
}
```

## Migration from Module

To migrate from the current `modules/configuration/gitops` module to the provider:

1. Replace the module call with the provider resource
2. Remove the module source reference
3. The provider handles all the complexity automatically

**Before (Module):**
```hcl
module "gitops" {
  source = "../../../../modules/configuration/gitops"

  cluster_id   = data.terraform_remote_state.infrastructure.outputs.cluster_id
  cluster_name = data.terraform_remote_state.infrastructure.outputs.cluster_name
  api_url      = data.terraform_remote_state.infrastructure.outputs.api_url

  labels = var.labels
  operator_version = "1.18.2"
}
```

**After (Provider):**
```hcl
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"
  version   = "1.18.2"

  labels = var.labels
}
```

## Resource Attributes

### Required
- `name` - Operator name (e.g., "openshift-gitops-operator")
- `namespace` - Namespace for installation
- `channel` - Operator channel (e.g., "latest", "stable")
- `source` - Catalog source (e.g., "redhat-operators")

### Optional
- `version` - Pin to specific version (auto-sets Manual approval)
- `install_plan_approval` - "Automatic" or "Manual" (default: "Automatic")
- `create_namespace` - Create namespace if missing (default: true)
- `namespace_labels` - Labels for namespace
- `operator_group_target_namespaces` - Target namespaces (empty = cluster-wide)
- `wait_for_csv` - Wait for CSV Succeeded phase (default: true)
- `wait_timeout` - Timeout for CSV wait (default: "10m")
- `labels` - Labels for operator resources

### Computed
- `installed_csv` - Name of installed CSV
- `csv_phase` - Current CSV phase

## Next Steps

- [ ] Test with real OpenShift cluster
- [ ] Add update support for version/channel changes
- [ ] Add support for operator uninstallation
- [ ] Publish to Terraform Registry
