# Troubleshooting Guide

This guide helps you troubleshoot common issues when using the OpenShift Operator Terraform Provider.

## Common Issues

### Operator Installation Fails

**Symptoms**: Resource creation fails with timeout or error messages.

**Possible Causes**:
1. Invalid operator name or channel
2. Catalog source not available
3. Insufficient permissions
4. Network issues

**Solutions**:
1. Verify the operator name matches the package name in the catalog:
   ```bash
   oc get packagemanifests -n openshift-marketplace | grep <operator-name>
   ```

2. Check available channels:
   ```bash
   oc get packagemanifests <operator-name> -n openshift-marketplace -o yaml | grep -A 10 channels
   ```

3. Verify catalog source exists:
   ```bash
   oc get catalogsource -n openshift-marketplace
   ```

4. Check RBAC permissions:
   ```bash
   oc auth can-i create subscriptions -n <namespace>
   oc auth can-i create operatorgroups -n <namespace>
   ```

### CSV Stuck in Installing Phase

**Symptoms**: `csv_phase` remains `Installing` for extended periods.

**Possible Causes**:
1. Operator dependencies not met
2. Resource constraints
3. Operator configuration issues

**Solutions**:
1. Check CSV status and conditions:
   ```bash
   oc get csv <csv-name> -n <namespace> -o yaml
   ```

2. Check for dependency issues:
   ```bash
   oc describe csv <csv-name> -n <namespace>
   ```

3. Review operator logs:
   ```bash
   oc logs -n <namespace> -l app=<operator-name>
   ```

### InstallPlan Approval Fails

**Symptoms**: Error when trying to approve InstallPlan manually.

**Possible Causes**:
1. Multiple InstallPlans exist
2. InstallPlan already approved
3. InstallPlan in failed state

**Solutions**:
1. List InstallPlans:
   ```bash
   oc get installplan -n <namespace>
   ```

2. Check InstallPlan status:
   ```bash
   oc get installplan <installplan-name> -n <namespace> -o yaml
   ```

3. If multiple InstallPlans exist, the provider will approve the one referenced by the Subscription.

### Import Fails

**Symptoms**: `terraform import` fails with error.

**Possible Causes**:
1. Incorrect resource ID format
2. Resource doesn't exist
3. Namespace/name mismatch

**Solutions**:
1. Verify resource ID format: `namespace/name`
   ```bash
   terraform import openshift_operator.example openshift-gitops-operator/openshift-gitops-operator
   ```

2. Verify Subscription exists:
   ```bash
   oc get subscription -n <namespace> <name>
   ```

3. Check that the namespace and name match exactly.

### Provider Authentication Fails

**Symptoms**: Provider cannot connect to cluster.

**Possible Causes**:
1. Invalid kubeconfig
2. Expired token
3. Network connectivity issues
4. TLS certificate issues

**Solutions**:
1. Test kubeconfig:
   ```bash
   kubectl --kubeconfig=<path> get nodes
   ```

2. Verify token is valid:
   ```bash
   oc whoami --token=<token>
   ```

3. For TLS issues, set `insecure = true` temporarily for testing:
   ```hcl
   provider "openshift" {
     insecure = true
   }
   ```

### Version Pinning Not Working

**Symptoms**: Operator upgrades despite version being pinned.

**Possible Causes**:
1. Version format incorrect
2. Version not available in channel
3. InstallPlan approval set to Automatic

**Solutions**:
1. Verify version format matches CSV name:
   ```bash
   oc get csv -n <namespace> | grep <operator-name>
   ```
   CSV names follow format: `<operator-name>.v<version>`

2. Check available versions in channel:
   ```bash
   oc get packagemanifests <operator-name> -n openshift-marketplace -o yaml
   ```

3. Ensure `install_plan_approval` is `Manual` when using version pinning.

## Debugging Tips

### Enable Debug Logging

Set the `TF_LOG` environment variable:

```bash
export TF_LOG=DEBUG
terraform apply
```

### Check Resource State

Inspect Terraform state:

```bash
terraform state show openshift_operator.example
```

### Verify Kubernetes Resources

Check all related resources:

```bash
# Subscription
oc get subscription -n <namespace> <name> -o yaml

# OperatorGroup
oc get operatorgroup -n <namespace> -o yaml

# InstallPlan
oc get installplan -n <namespace> -o yaml

# CSV
oc get csv -n <namespace> -o yaml
```

### Review Provider Logs

The provider logs detailed information about operations. Check logs for:
- Resource creation steps
- Wait conditions
- Error details

## Getting Help

If you're still experiencing issues:

1. Check the [GitHub Issues](https://github.com/rh-mobb/terraform-provider-openshift/issues)
2. Create a new issue with:
   - Provider version
   - Terraform version
   - OpenShift version
   - Terraform configuration (sanitized)
   - Error messages and logs
   - Steps to reproduce
