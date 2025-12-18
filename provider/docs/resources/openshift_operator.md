# openshift_operator Resource

Manages an OpenShift operator installation using the Operator Lifecycle Manager (OLM).

## Example Usage

### Basic Installation

```hcl
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

### With Custom Configuration

```hcl
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
  channel   = "latest"
  source    = "redhat-operators"

  install_plan_approval = "Manual"
  create_namespace     = true

  namespace_labels = {
    "openshift.io/cluster-monitoring" = "true"
  }

  operator_group_target_namespaces = []  # Empty = cluster-wide

  wait_for_csv = true
  wait_timeout = "15m"

  labels = {
    "app.kubernetes.io/managed-by" = "Terraform"
    "environment"                 = "production"
  }
}
```

## Argument Reference

### Required Arguments

- `name` (string) - Name of the operator (e.g., `openshift-gitops-operator`). This must match the operator's package name in the catalog.

- `namespace` (string) - Namespace where the operator will be installed. The provider will create this namespace if `create_namespace` is `true`.

- `channel` (string) - Channel for the operator subscription (e.g., `latest`, `stable`, `v1.18`). Available channels depend on the operator.

- `source` (string) - Catalog source name (e.g., `redhat-operators`, `certified-operators`, `community-operators`). The catalog source must exist in the `openshift-marketplace` namespace.

### Optional Arguments

- `version` (string) - Specific version to install (e.g., `1.18.2`). If specified, `install_plan_approval` is automatically set to `Manual`. The version must be available in the specified channel.

- `install_plan_approval` (string) - Install plan approval strategy. Valid values are `Automatic` or `Manual`. Defaults to `Automatic`. If `version` is specified, this is automatically set to `Manual`.

- `create_namespace` (bool) - Whether to create the namespace if it doesn't exist. Defaults to `true`.

- `namespace_labels` (map of strings) - Labels to apply to the namespace when it is created. Only applied if `create_namespace` is `true`.

- `operator_group_target_namespaces` (list of strings) - Target namespaces for the OperatorGroup. An empty list means the operator will be installed cluster-wide. Defaults to an empty list (cluster-wide).

- `wait_for_csv` (bool) - Whether to wait for the ClusterServiceVersion (CSV) to reach the `Succeeded` phase before considering the resource created. Defaults to `true`.

- `wait_timeout` (string) - Timeout for waiting for CSV to succeed. Accepts duration strings like `10m`, `1h`, `30s`. Defaults to `10m`.

- `labels` (map of strings) - Labels to apply to operator resources (Subscription, OperatorGroup). The provider automatically adds `app.kubernetes.io/managed-by: Terraform` and `app.kubernetes.io/name: <operator-name>` labels.

### Computed Attributes

- `installed_csv` (string) - Name of the installed ClusterServiceVersion (CSV). This is set after the CSV is created by OLM.

- `csv_phase` (string) - Current phase of the CSV. Common values include `Succeeded`, `Installing`, `Failed`, `Pending`. This is only populated if `wait_for_csv` is `true` or after the CSV is created.

## Timeouts

The resource supports the following timeout configurations:

- `create` - Default: 20 minutes
- `update` - Default: 20 minutes
- `delete` - Default: 10 minutes

Example:

```hcl
resource "openshift_operator" "gitops" {
  # ... other configuration ...

  timeouts {
    create = "30m"
    update = "30m"
    delete = "15m"
  }
}
```

## Import

Existing operator installations can be imported using the resource ID format: `namespace/name`

```bash
terraform import openshift_operator.gitops openshift-gitops-operator/openshift-gitops-operator
```

## Lifecycle Management

The provider handles the following resources automatically:

1. **Namespace** - Created if `create_namespace` is `true` and the namespace doesn't exist
2. **OperatorGroup** - Created with the specified target namespaces
3. **Subscription** - Created with the specified channel, source, and approval strategy
4. **InstallPlan** - Automatically approved if `install_plan_approval` is `Manual`
5. **ClusterServiceVersion (CSV)** - Monitored until it reaches `Succeeded` phase (if `wait_for_csv` is `true`)

## Version Pinning

When `version` is specified:

- The `startingCSV` field is set in the Subscription
- `install_plan_approval` is automatically set to `Manual`
- The provider will approve the InstallPlan automatically
- The operator will be pinned to the specified version and won't auto-upgrade

## Update Behavior

The provider supports updating the following attributes:

- `channel` - Changes the subscription channel
- `source` - Changes the catalog source
- `install_plan_approval` - Changes the approval strategy
- `version` - Updates to a new version (requires Manual approval)
- `labels` - Updates resource labels

When updating `version` or changing `install_plan_approval` to `Manual`, the provider will:

1. Update the Subscription
2. Wait for a new InstallPlan
3. Automatically approve the InstallPlan
4. Wait for the CSV to be updated

## Notes

- The provider does not delete the namespace or CSV when the resource is destroyed. The CSV will be cleaned up by OLM when the Subscription is deleted.
- If an operator installation fails, check the CSV status and InstallPlan for error details.
- Some operators may require additional configuration after installation (e.g., CustomResource instances).
- The provider requires appropriate RBAC permissions to create and manage operator resources.

## See Also

- [OpenShift Operator Lifecycle Manager Documentation](https://docs.openshift.com/container-platform/latest/operators/understanding/olm/olm-understanding-olm.html)
- [OLM Documentation](https://olm.operatorframework.io/)
