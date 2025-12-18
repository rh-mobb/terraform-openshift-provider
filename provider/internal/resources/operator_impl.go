package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/redhat/terraform-provider-openshift-operator/internal/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
)

func resourceOperatorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ok := meta.(*client.Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("invalid client type"))
	}
	dynamicClient := c.Dynamic

	namespace, ok := d.Get("namespace").(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("namespace must be a string"))
	}
	name, ok := d.Get("name").(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("name must be a string"))
	}
	channel, ok := d.Get("channel").(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("channel must be a string"))
	}
	source, ok := d.Get("source").(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("source must be a string"))
	}
	version, ok := d.Get("version").(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("version must be a string"))
	}
	installPlanApproval, ok := d.Get("install_plan_approval").(string)
	if !ok {
		return diag.FromErr(fmt.Errorf("install_plan_approval must be a string"))
	}
	shouldCreateNamespace, ok := d.Get("create_namespace").(bool)
	if !ok {
		return diag.FromErr(fmt.Errorf("create_namespace must be a bool"))
	}
	waitForCSV, ok := d.Get("wait_for_csv").(bool)
	if !ok {
		return diag.FromErr(fmt.Errorf("wait_for_csv must be a bool"))
	}

	// Validate inputs
	if namespace == "" {
		return diag.FromErr(fmt.Errorf("namespace cannot be empty"))
	}
	if name == "" {
		return diag.FromErr(fmt.Errorf("name cannot be empty"))
	}
	if channel == "" {
		return diag.FromErr(fmt.Errorf("channel cannot be empty"))
	}
	if source == "" {
		return diag.FromErr(fmt.Errorf("source cannot be empty"))
	}

	// If version is specified, force Manual approval
	if version != "" {
		installPlanApproval = "Manual"
	}

	// Build labels
	labels := make(map[string]interface{})
	if v, ok := d.GetOk("labels"); ok {
		if labelMap, ok := v.(map[string]interface{}); ok {
			for k, v := range labelMap {
				labels[k] = v
			}
		}
	}
	labels["app.kubernetes.io/managed-by"] = "Terraform"
	labels["app.kubernetes.io/name"] = name

	// Build namespace labels
	namespaceLabels := make(map[string]string)
	if v, ok := d.GetOk("namespace_labels"); ok {
		if labelMap, ok := v.(map[string]interface{}); ok {
			for k, v := range labelMap {
				if strVal, ok := v.(string); ok {
					namespaceLabels[k] = strVal
				}
			}
		}
	}

	// 1. Create namespace if needed
	if shouldCreateNamespace {
		if err := createNamespace(ctx, dynamicClient, namespace, namespaceLabels); err != nil {
			// Check if namespace already exists
			if !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "AlreadyExists") {
				return diag.FromErr(fmt.Errorf("failed to create namespace %s: %w", namespace, err))
			}
		}
	}

	// 2. Create OperatorGroup
	operatorGroupName := fmt.Sprintf("%s-operatorgroup", name)
	targetNamespaces, ok := d.Get("operator_group_target_namespaces").([]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("operator_group_target_namespaces must be a list"))
	}
	if err := createOperatorGroup(ctx, dynamicClient, namespace, operatorGroupName, targetNamespaces, labels); err != nil {
		if !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "AlreadyExists") {
			return diag.FromErr(fmt.Errorf("failed to create OperatorGroup %s in namespace %s: %w", operatorGroupName, namespace, err))
		}
	}

	// 3. Create Subscription
	startingCSV := ""
	if version != "" {
		// CSV name format: {name}.v{version}
		startingCSV = fmt.Sprintf("%s.v%s", name, version)
	}

	sourceNamespace := "openshift-marketplace"
	if err := createSubscription(ctx, dynamicClient, namespace, name, channel, source, sourceNamespace, installPlanApproval, startingCSV, labels); err != nil {
		if !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "AlreadyExists") {
			return diag.FromErr(fmt.Errorf("failed to create Subscription %s in namespace %s: %w", name, namespace, err))
		}
	}

	// 4. Handle Manual approval if needed
	var csvName string
	if installPlanApproval == "Manual" {
		// Wait for InstallPlan reference
		timeout := 5 * time.Minute
		installPlanName, err := waitForInstallPlanRef(ctx, dynamicClient, namespace, name, timeout)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to get InstallPlan reference",
					Detail:   fmt.Sprintf("Operator: %s, Namespace: %s. Error: %v. This may indicate the operator is not available in the specified channel or catalog source.", name, namespace, err),
				},
			}
		}

		// Approve InstallPlan
		if err := approveInstallPlan(ctx, dynamicClient, namespace, installPlanName); err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to approve InstallPlan",
					Detail:   fmt.Sprintf("InstallPlan: %s, Namespace: %s. Error: %v. Check InstallPlan status for details.", installPlanName, namespace, err),
				},
			}
		}

		// Wait for installedCSV
		timeout = 10 * time.Minute
		csvName, err = waitForInstalledCSV(ctx, dynamicClient, namespace, name, timeout)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to wait for installed CSV",
					Detail:   fmt.Sprintf("Operator: %s, Namespace: %s. Error: %v. The InstallPlan may have failed or the operator installation is taking longer than expected.", name, namespace, err),
				},
			}
		}
	} else {
		// Automatic approval - wait for installedCSV
		timeout := 10 * time.Minute
		var err error
		csvName, err = waitForInstalledCSV(ctx, dynamicClient, namespace, name, timeout)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Failed to wait for installed CSV",
					Detail:   fmt.Sprintf("Operator: %s, Namespace: %s. Error: %v. Check Subscription and InstallPlan status for details.", name, namespace, err),
				},
			}
		}
	}

	// 5. Wait for CSV to be Succeeded if requested
	var csvPhase string
	if waitForCSV {
		waitTimeoutStr, ok := d.Get("wait_timeout").(string)
		if !ok {
			return diag.FromErr(fmt.Errorf("wait_timeout must be a string"))
		}
		waitTimeout, err := time.ParseDuration(waitTimeoutStr)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid wait_timeout format: %w", err))
		}

		csvPhase, err = waitForCSVSucceeded(ctx, dynamicClient, namespace, csvName, waitTimeout)
		if err != nil {
			// Get CSV details for better error message
			csv, csvErr := getCSV(ctx, dynamicClient, namespace, csvName)
			detail := fmt.Sprintf("CSV %s in namespace %s did not reach Succeeded phase within timeout. Current phase: %s", csvName, namespace, csvPhase)
			if csvErr == nil {
				if _, found, _ := unstructured.NestedSlice(csv.Object, "status", "conditions"); found {
					detail += ". Check CSV conditions for more details."
				}
			}
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "CSV did not reach Succeeded phase",
					Detail:   detail,
				},
			}
		}
	} else {
		// Still read the phase even if not waiting
		csv, err := getCSV(ctx, dynamicClient, namespace, csvName)
		if err == nil {
			if status, found, _ := unstructured.NestedMap(csv.Object, "status"); found {
				if phase, found, _ := unstructured.NestedString(status, "phase"); found {
					csvPhase = phase
				}
			}
		}
	}

	// Set ID and state
	d.SetId(fmt.Sprintf("%s/%s", namespace, name))
	_ = d.Set("installed_csv", csvName)
	_ = d.Set("csv_phase", csvPhase)

	return resourceOperatorRead(ctx, d, meta)
}

func resourceOperatorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ok := meta.(*client.Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("invalid client type"))
	}
	dynamicClient := c.Dynamic

	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 2 {
		return diag.FromErr(fmt.Errorf("invalid resource ID format: %s", d.Id()))
	}

	namespace := idParts[0]
	name := idParts[1]

	// Read Subscription
	sub, err := getSubscription(ctx, dynamicClient, namespace, name)
	if err != nil {
		d.SetId("")
		return nil
	}

	// Read Subscription spec to populate state
	if spec, found, _ := unstructured.NestedMap(sub.Object, "spec"); found {
		if channel, found, _ := unstructured.NestedString(spec, "channel"); found {
			_ = d.Set("channel", channel)
		}
		if source, found, _ := unstructured.NestedString(spec, "source"); found {
			_ = d.Set("source", source)
		}
		if installPlanApproval, found, _ := unstructured.NestedString(spec, "installPlanApproval"); found {
			_ = d.Set("install_plan_approval", installPlanApproval)
		}
		if startingCSV, found, _ := unstructured.NestedString(spec, "startingCSV"); found && startingCSV != "" {
			// Extract version from CSV name format: {name}.v{version}
			// This is a best-effort extraction
			if strings.HasPrefix(startingCSV, name+".v") {
				version := strings.TrimPrefix(startingCSV, name+".v")
				_ = d.Set("version", version)
			}
		}
	}

	// Read metadata
	if metadata, found, _ := unstructured.NestedMap(sub.Object, "metadata"); found {
		if labels, found, _ := unstructured.NestedStringMap(metadata, "labels"); found {
			// Filter out Terraform-managed labels
			resourceLabels := make(map[string]interface{})
			for k, v := range labels {
				if k != "app.kubernetes.io/managed-by" && k != "app.kubernetes.io/name" {
					resourceLabels[k] = v
				}
			}
			if len(resourceLabels) > 0 {
				_ = d.Set("labels", resourceLabels)
			}
		}
	}

	// Read CSV if available
	if status, found, _ := unstructured.NestedMap(sub.Object, "status"); found {
		if csvName, found, _ := unstructured.NestedString(status, "installedCSV"); found && csvName != "" {
			_ = d.Set("installed_csv", csvName)

			// Read CSV phase
			csv, err := getCSV(ctx, dynamicClient, namespace, csvName)
			if err == nil {
				if csvStatus, found, _ := unstructured.NestedMap(csv.Object, "status"); found {
					if phase, found, _ := unstructured.NestedString(csvStatus, "phase"); found {
						_ = d.Set("csv_phase", phase)
					}
				}
			}
		}
	}

	// Set namespace and name
	_ = d.Set("namespace", namespace)
	_ = d.Set("name", name)

	return nil
}

func resourceOperatorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ok := meta.(*client.Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("invalid client type"))
	}
	dynamicClient := c.Dynamic

	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 2 {
		return diag.FromErr(fmt.Errorf("invalid resource ID format: %s", d.Id()))
	}

	namespace := idParts[0]
	name := idParts[1]

	// Get Subscription
	sub, err := getSubscription(ctx, dynamicClient, namespace, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to get Subscription: %w", err))
	}

	needsUpdate := false

	// Update channel if changed
	if d.HasChange("channel") {
		channel, ok := d.Get("channel").(string)
		if !ok {
			return diag.FromErr(fmt.Errorf("channel must be a string"))
		}
		if err := unstructured.SetNestedField(sub.Object, channel, "spec", "channel"); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update channel: %w", err))
		}
		needsUpdate = true
	}

	// Update source if changed
	if d.HasChange("source") {
		source, ok := d.Get("source").(string)
		if !ok {
			return diag.FromErr(fmt.Errorf("source must be a string"))
		}
		if err := unstructured.SetNestedField(sub.Object, source, "spec", "source"); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update source: %w", err))
		}
		needsUpdate = true
	}

	// Update install plan approval if changed
	if d.HasChange("install_plan_approval") {
		installPlanApproval, ok := d.Get("install_plan_approval").(string)
		if !ok {
			return diag.FromErr(fmt.Errorf("install_plan_approval must be a string"))
		}
		if err := unstructured.SetNestedField(sub.Object, installPlanApproval, "spec", "installPlanApproval"); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update installPlanApproval: %w", err))
		}
		needsUpdate = true
	}

	// Update version/startingCSV if changed
	if d.HasChange("version") {
		version, ok := d.Get("version").(string)
		if !ok {
			return diag.FromErr(fmt.Errorf("version must be a string"))
		}
		if version != "" {
			// If version is specified, force Manual approval
			if err := unstructured.SetNestedField(sub.Object, "Manual", "spec", "installPlanApproval"); err != nil {
				return diag.FromErr(fmt.Errorf("failed to set installPlanApproval to Manual: %w", err))
			}
			startingCSV := fmt.Sprintf("%s.v%s", name, version)
			if err := unstructured.SetNestedField(sub.Object, startingCSV, "spec", "startingCSV"); err != nil {
				return diag.FromErr(fmt.Errorf("failed to update startingCSV: %w", err))
			}
		} else {
			// Remove startingCSV if version is cleared
			unstructured.RemoveNestedField(sub.Object, "spec", "startingCSV")
		}
		needsUpdate = true
	}

	// Update labels if changed
	if d.HasChange("labels") {
		labels := make(map[string]interface{})
		if v, ok := d.GetOk("labels"); ok {
			if labelMap, ok := v.(map[string]interface{}); ok {
				for k, v := range labelMap {
					labels[k] = v
				}
			}
		}
		labels["app.kubernetes.io/managed-by"] = "Terraform"
		labels["app.kubernetes.io/name"] = name

		if err := unstructured.SetNestedMap(sub.Object, labels, "metadata", "labels"); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update labels: %w", err))
		}
		needsUpdate = true
	}

	// Apply updates if needed
	if needsUpdate {
		subGVR := k8sschema.GroupVersionResource{
			Group:    "operators.coreos.com",
			Version:  "v1alpha1",
			Resource: "subscriptions",
		}

		_, err := dynamicClient.Resource(subGVR).Namespace(namespace).Update(ctx, sub, metav1.UpdateOptions{})
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Subscription: %w", err))
		}

		// If Manual approval and version changed, approve new InstallPlan
		installPlanApproval, ok := d.Get("install_plan_approval").(string)
		if !ok {
			return diag.FromErr(fmt.Errorf("install_plan_approval must be a string"))
		}
		if installPlanApproval == "Manual" || d.HasChange("version") {
			// Wait for new InstallPlan reference
			timeout := 5 * time.Minute
			installPlanName, err := waitForInstallPlanRef(ctx, dynamicClient, namespace, name, timeout)
			if err == nil {
				// Approve InstallPlan
				if err := approveInstallPlan(ctx, dynamicClient, namespace, installPlanName); err != nil {
					return diag.FromErr(fmt.Errorf("failed to approve InstallPlan: %w", err))
				}

				// Wait for installedCSV update
				timeout = 10 * time.Minute
				csvName, err := waitForInstalledCSV(ctx, dynamicClient, namespace, name, timeout)
				if err == nil {
					_ = d.Set("installed_csv", csvName)

					// Wait for CSV to succeed if requested
					waitForCSV, ok := d.Get("wait_for_csv").(bool)
					if !ok {
						return diag.FromErr(fmt.Errorf("wait_for_csv must be a bool"))
					}
					if waitForCSV {
						waitTimeoutStr, ok := d.Get("wait_timeout").(string)
						if !ok {
							return diag.FromErr(fmt.Errorf("wait_timeout must be a string"))
						}
						waitTimeout, err := time.ParseDuration(waitTimeoutStr)
						if err == nil {
							csvPhase, err := waitForCSVSucceeded(ctx, dynamicClient, namespace, csvName, waitTimeout)
							if err == nil {
								_ = d.Set("csv_phase", csvPhase)
							}
						}
					}
				}
			}
		}
	}

	return resourceOperatorRead(ctx, d, meta)
}

func resourceOperatorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, ok := meta.(*client.Client)
	if !ok {
		return diag.FromErr(fmt.Errorf("invalid client type"))
	}
	dynamicClient := c.Dynamic

	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 2 {
		return diag.FromErr(fmt.Errorf("invalid resource ID format: %s", d.Id()))
	}

	namespace := idParts[0]
	name := idParts[1]

	// Delete Subscription
	subGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "subscriptions",
	}

	err := dynamicClient.Resource(subGVR).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return diag.FromErr(fmt.Errorf("failed to delete Subscription: %w", err))
	}

	// Delete OperatorGroup
	ogGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1",
		Resource: "operatorgroups",
	}

	operatorGroupName := fmt.Sprintf("%s-operatorgroup", name)
	err = dynamicClient.Resource(ogGVR).Namespace(namespace).Delete(ctx, operatorGroupName, metav1.DeleteOptions{})
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return diag.FromErr(fmt.Errorf("failed to delete OperatorGroup: %w", err))
	}

	// Note: We don't delete the namespace or CSV as they may be managed separately
	// The CSV will be cleaned up by OLM when the Subscription is deleted

	return nil
}
