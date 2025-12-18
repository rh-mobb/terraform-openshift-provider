---
subcategory: ""
layout: "openshift"
page_title: "openshift_operator Data Source"
description: |-
  Reads the status and configuration of an existing OpenShift operator installation.
---

# openshift_operator

Reads the status and configuration of an existing OpenShift operator installation.

## Example Usage

```hcl
data "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
}

output "gitops_version" {
  value = data.openshift_operator.gitops.installed_csv_version
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the operator.
* `namespace` - (Required) Namespace where the operator is installed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `channel` - Current channel for the operator subscription.
* `source` - Catalog source name.
* `version` - Pinned version if version pinning is enabled.
* `install_plan_approval` - Install plan approval strategy.
* `installed_csv` - Name of the installed CSV.
* `csv_phase` - Current phase of the CSV.
* `csv_version` - Version of the installed CSV.
* `subscription_state` - State of the subscription.
* `current_csv` - Current CSV being tracked.
* `installed_csv_version` - Version extracted from CSV name.
