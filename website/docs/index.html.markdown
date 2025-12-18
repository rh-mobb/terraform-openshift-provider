---
subcategory: ""
layout: "openshift"
page_title: "OpenShift Provider"
description: |-
  The OpenShift provider is used to manage OpenShift operators and related resources.
---

# OpenShift Provider

The OpenShift provider is used to manage OpenShift operators and related resources using the Operator Lifecycle Manager (OLM).

Use the navigation to the left to read about the available resources and data sources.

## Example Usage

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

## Authentication

The OpenShift provider supports multiple authentication methods:

1. **Kubeconfig file** (recommended for local development)
2. **Host and token** (useful for CI/CD)
3. **In-cluster config** (when running inside Kubernetes)
4. **Default kubeconfig locations** (automatically detected)

See the [Provider Configuration](https://registry.terraform.io/providers/rh-mobb/openshift/latest/docs) documentation for details.
