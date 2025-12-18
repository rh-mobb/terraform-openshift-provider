---
subcategory: ""
layout: "openshift"
page_title: "openshift_operator Resource"
description: |-
  Manages an OpenShift operator installation using the Operator Lifecycle Manager (OLM).
---

# openshift_operator

Manages an OpenShift operator installation using the Operator Lifecycle Manager (OLM).

## Example Usage

```hcl
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the operator.
* `namespace` - (Required) Namespace where the operator will be installed.
* `channel` - (Required) Channel for the operator subscription.
* `source` - (Required) Catalog source name.
* `version` - (Optional) Specific version to install.
* `install_plan_approval` - (Optional) Install plan approval strategy. Defaults to `Automatic`.
* `create_namespace` - (Optional) Whether to create the namespace. Defaults to `true`.
* `namespace_labels` - (Optional) Labels to apply to the namespace.
* `operator_group_target_namespaces` - (Optional) Target namespaces for the OperatorGroup.
* `wait_for_csv` - (Optional) Whether to wait for CSV to succeed. Defaults to `true`.
* `wait_timeout` - (Optional) Timeout for waiting for CSV. Defaults to `10m`.
* `labels` - (Optional) Labels to apply to operator resources.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `installed_csv` - Name of the installed ClusterServiceVersion (CSV).
* `csv_phase` - Current phase of the CSV.

## Import

Operator installations can be imported using the resource ID format: `namespace/name`

```bash
terraform import openshift_operator.gitops openshift-gitops-operator/openshift-gitops-operator
```
