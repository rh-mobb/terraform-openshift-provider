# OpenShift Operator Terraform Provider - Design Document

## Current Status

The provider has been extracted from the original Terraform module repository into its own dedicated repository (`rh-mobb_openshift-provider`) to prepare for production release and publication to the Terraform Registry.

**Next Steps:**
- Make the provider production-ready
- Add comprehensive tests
- Complete documentation
- Publish to GitHub under the `rh-mobb` organization
- Publish to Terraform Registry under the `rh-mobb` namespace

## Problem Statement

Installing OpenShift operators via Terraform using the Kubernetes provider is complex and error-prone:

- Requires managing multiple resources (Namespace, OperatorGroup, Subscription, InstallPlan, CSV)
- Complex wait logic that differs between Manual and Automatic approval
- Need to handle InstallPlan approval manually
- Status fields not directly accessible from `kubernetes_manifest` resources
- Multiple `null_resource` workarounds with `local-exec` provisioners

## Proposed Solution

Create a custom Terraform provider `terraform-provider-openshift` that encapsulates all operator lifecycle management into a single resource, with the ability to add other OpenShift-related resources in the future.

## Resource Design

```hcl
resource "openshift_operator" "gitops" {
  name      = "openshift-gitops-operator"
  namespace = "openshift-gitops-operator"

  # Operator source configuration
  channel   = "latest"              # or "stable", "v1.18", etc.
  source    = "redhat-operators"    # Catalog source name
  version   = "1.18.2"              # Optional: pin to specific version

  # Approval strategy
  install_plan_approval = "Manual"  # or "Automatic"

  # Namespace configuration
  create_namespace = true
  namespace_labels = {
    "openshift.io/cluster-monitoring" = "true"
  }

  # OperatorGroup configuration
  operator_group_target_namespaces = []  # Empty = cluster-wide

  # Wait configuration
  wait_for_csv = true
  wait_timeout = "10m"

  # Labels
  labels = {
    "app.kubernetes.io/managed-by" = "Terraform"
  }
}
```

## Provider Features

### 1. Lifecycle Management
- **Create**: Creates namespace (if needed), OperatorGroup, Subscription, approves InstallPlan (if Manual), waits for CSV
- **Read**: Reads current state of operator installation
- **Update**: Handles version upgrades, channel changes, approval strategy changes
- **Delete**: Removes Subscription, OperatorGroup, optionally namespace

### 2. Automatic InstallPlan Approval
- When `install_plan_approval = "Manual"` and `version` is specified:
  - Automatically finds and approves the InstallPlan
  - Waits for CSV installation to complete
  - Handles edge cases (multiple InstallPlans, failed InstallPlans)

### 3. Version Pinning
- When `version` is specified:
  - Sets `startingCSV` in Subscription
  - Sets `installPlanApproval = "Manual"` automatically
  - Prevents automatic upgrades

### 4. State Management
- Tracks operator installation state
- Handles partial failures gracefully
- Provides detailed error messages

## Implementation Approach

### Option 1: Go-based Provider (Recommended)
- Use Terraform Plugin SDK v2
- Use Kubernetes Go client libraries
- Implement proper CRUD operations
- Handle retries and timeouts internally

### Option 2: gRPC-based Provider
- Use Terraform Plugin Framework
- More modern but requires more setup

## Benefits

1. **Simplicity**: Single resource instead of 5+ resources
2. **Reliability**: Handles edge cases internally
3. **Maintainability**: Centralized logic for operator management
4. **Reusability**: Can be used across multiple projects
5. **Better Error Messages**: Provider-specific error handling

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
  version   = "1.18.2"  # Automatically sets Manual approval
}
```

### Multiple Operators
```hcl
resource "openshift_operator" "gitops" {
  name    = "openshift-gitops-operator"
  channel = "latest"
  source  = "redhat-operators"
  version = "1.18.2"
}

resource "openshift_operator" "service_mesh" {
  name    = "servicemeshoperator"
  channel = "stable"
  source  = "redhat-operators"
}
```

## Production Readiness Checklist

### Code Quality
- [ ] Complete implementation of all CRUD operations
- [ ] Comprehensive error handling and user-friendly error messages
- [ ] Input validation for all resource attributes
- [ ] Proper logging and debugging support
- [ ] Code comments and documentation
- [ ] Follow Go best practices and conventions

### Testing
- [ ] Unit tests for core functionality
- [ ] Integration tests against real OpenShift cluster
- [ ] Acceptance tests using Terraform's test framework
- [ ] Test coverage for edge cases (failed InstallPlans, multiple InstallPlans, etc.)
- [ ] Test matrix for different OpenShift versions
- [ ] CI/CD pipeline for automated testing

### Documentation
- [ ] Provider README with installation instructions
- [ ] Resource documentation with all attributes explained
- [ ] Usage examples for common scenarios
- [ ] Migration guide from module-based approach
- [ ] Troubleshooting guide
- [ ] CHANGELOG.md for version tracking
- [ ] LICENSE file

### Terraform Registry Requirements
- [ ] Provider metadata file (`metadata.json`)
- [ ] GPG signing key setup
- [ ] Release process documentation
- [ ] Version tagging strategy (semantic versioning)
- [ ] GitHub releases with binaries for multiple platforms

### GitHub Publication
- [ ] Repository setup under `rh-mobb` organization
- [ ] Repository description and topics
- [ ] GitHub Actions workflows for CI/CD
- [ ] Issue and PR templates
- [ ] Contributing guidelines
- [ ] Code of conduct (if applicable)
- [ ] Security policy

### Terraform Registry Publication
- [ ] Create provider namespace: `rh-mobb/openshift`
- [ ] Configure GPG key in Terraform Registry
- [ ] Set up GitHub release automation
- [ ] Initial release (v0.1.0 or v1.0.0)
- [ ] Registry documentation review and approval

## Testing Strategy

### Unit Tests
- Test resource schema validation
- Test state conversion logic
- Test InstallPlan approval logic
- Test version pinning logic
- Mock Kubernetes API responses

### Integration Tests
- Test against real OpenShift cluster (or kind/k3d)
- Test operator installation lifecycle
- Test upgrade scenarios
- Test deletion scenarios
- Test error handling (invalid channels, missing operators, etc.)

### Acceptance Tests
- Use Terraform's acceptance test framework
- Test full resource lifecycle
- Test state persistence
- Test concurrent resource creation

### Test Infrastructure
- Use GitHub Actions for CI/CD
- Support testing against multiple OpenShift versions
- Use test fixtures for reproducible tests
- Clean up test resources properly

## Documentation Requirements

### Provider Documentation
- **README.md**: Overview, installation, quick start
- **docs/**: Detailed documentation directory
  - `resources/openshift_operator.md`: Complete resource reference
  - `guides/`: Usage guides and examples
  - `troubleshooting.md`: Common issues and solutions

### Resource Documentation
- All attributes with descriptions
- Required vs optional fields
- Default values
- Examples for each use case
- Import instructions
- Output attributes

### Examples
- Basic operator installation
- Version-pinned installation
- Multiple operators
- Custom namespace configuration
- OperatorGroup configuration
- Upgrade scenarios

## Publication Plan

### GitHub Repository
- **Organization**: `rh-mobb`
- **Repository**: `terraform-provider-openshift`
- **Visibility**: Public
- **License**: Apache 2.0 (or as specified)

### Terraform Registry
- **Namespace**: `rh-mobb`
- **Provider Name**: `openshift`
- **Full Path**: `registry.terraform.io/rh-mobb/openshift`
- **Initial Version**: v0.1.0 (or v1.0.0 if stable)

### Release Process
1. Tag releases using semantic versioning (v0.1.0, v0.2.0, etc.)
2. Create GitHub releases with release notes
3. Build and attach binaries for:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)
4. Terraform Registry will automatically detect and publish releases

## Open Questions

1. Should the provider handle OperatorGroup creation automatically or require it separately?
2. Should namespace creation be optional or always handled?
3. How to handle operator upgrades when version changes?
4. Should we support uninstalling operators (removing CSV)?
5. What OpenShift versions should we support? (minimum version)
6. Should we support multiple catalog sources?
7. How to handle operator dependencies?

## References

### Terraform Provider Development
- [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk)
- [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
- [Terraform Provider Development Guide](https://developer.hashicorp.com/terraform/plugin)
- [Terraform Registry Publishing Guide](https://developer.hashicorp.com/terraform/registry/providers/publishing)

### Kubernetes/OpenShift
- [Kubernetes Go Client](https://github.com/kubernetes/client-go)
- [OLM Documentation](https://olm.operatorframework.io/)
- [OpenShift Operator Lifecycle Manager](https://docs.openshift.com/container-platform/latest/operators/understanding/olm/olm-understanding-olm.html)

### Testing
- [Terraform Acceptance Testing](https://www.terraform.io/plugin/sdkv2/testing/acceptance-tests)
- [Go Testing Best Practices](https://golang.org/doc/effective_go#testing)
