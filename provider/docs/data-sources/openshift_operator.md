# openshift_operator Data Source

Reads the status and configuration of an existing OpenShift operator installation.

This data source is useful for:
- Checking if an operator is installed
- Reading operator version and status information
- Using operator status in conditional logic
- Discovering existing operator configurations

## Example Usage

### Basic Usage

```hcl
data "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
}

output "gitops_version" {
  value = data.openshift_operator.gitops.installed_csv_version
}

output "gitops_status" {
  value = data.openshift_operator.gitops.csv_phase
}
```

### Conditional Logic

```hcl
data "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
}

# Only create a resource if operator is installed and succeeded
resource "some_resource" "example" {
  count = data.openshift_operator.gitops.csv_phase == "Succeeded" ? 1 : 0
  # ...
}
```

### Version Check

```hcl
data "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"
}

# Check if operator version meets requirements
locals {
  gitops_version_ok = data.openshift_operator.gitops.csv_version >= "1.18.0"
}
```

## Argument Reference

### Required Arguments

- `name` (string) - Name of the operator (e.g., `openshift-gitops-operator`). This must match the operator's Subscription name.

- `namespace` (string) - Namespace where the operator is installed.

## Attributes Reference

All attributes are computed (read-only).

- `channel` (string) - Current channel for the operator subscription (e.g., `latest`, `stable`).

- `source` (string) - Catalog source name for the operator subscription (e.g., `redhat-operators`).

- `version` (string) - Pinned version (extracted from `startingCSV`) if version pinning is enabled. Empty if not pinned.

- `install_plan_approval` (string) - Install plan approval strategy. Values: `Automatic` or `Manual`.

- `installed_csv` (string) - Name of the installed ClusterServiceVersion (CSV).

- `csv_phase` (string) - Current phase of the CSV. Common values:
  - `Succeeded` - CSV is installed and ready
  - `Installing` - CSV is being installed
  - `Failed` - CSV installation failed
  - `Pending` - CSV is pending installation
  - `Deleting` - CSV is being deleted

- `csv_version` (string) - Version of the installed CSV (from CSV spec).

- `installed_csv_version` (string) - Version extracted from the installed CSV name (format: `{name}.v{version}`).

- `subscription_state` (string) - State of the subscription. Common values:
  - `AtLatestKnown` - Subscription is at the latest known version
  - `UpgradePending` - An upgrade is pending
  - `UpgradeAvailable` - An upgrade is available

- `current_csv` (string) - Current CSV that the subscription is tracking (may differ from `installed_csv` during upgrades).

## Notes

- The data source will return an error if the operator Subscription is not found.
- If the CSV is not yet installed, `installed_csv` and related CSV attributes will be empty.
- The data source does not wait for the operator to be ready - it reads the current state immediately.
- Use the `csv_phase` attribute to check if the operator is ready before using it in other resources.

## See Also

- [openshift_operator Resource](../resources/openshift_operator.md) - For managing operator installations
