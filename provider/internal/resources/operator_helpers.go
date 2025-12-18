package resources

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8sschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
)

const (
	subscriptionGVR = "operators.coreos.com/v1alpha1"
	subscriptionKind = "Subscription"
	operatorGroupGVR = "operators.coreos.com/v1"
	operatorGroupKind = "OperatorGroup"
	installPlanGVR   = "operators.coreos.com/v1alpha1"
	installPlanKind  = "InstallPlan"
	csvGVR           = "operators.coreos.com/v1alpha1"
	csvKind          = "ClusterServiceVersion"
)

// Helper functions for operator resource management

func createNamespace(ctx context.Context, client dynamic.Interface, namespace string, labels map[string]string) error {
	nsGVR := k8sschema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}

	ns := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Namespace",
			"metadata": map[string]interface{}{
				"name":   namespace,
				"labels": labels,
			},
		},
	}

	_, err := client.Resource(nsGVR).Create(ctx, ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace %s: %w", namespace, err)
	}

	return nil
}

func createOperatorGroup(ctx context.Context, client dynamic.Interface, namespace, name string, targetNamespaces []interface{}, labels map[string]interface{}) error {
	ogGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1",
		Resource: "operatorgroups",
	}

	og := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "operators.coreos.com/v1",
			"kind":       "OperatorGroup",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
				"labels":    labels,
			},
			"spec": map[string]interface{}{
				"targetNamespaces": targetNamespaces,
			},
		},
	}

	_, err := client.Resource(ogGVR).Namespace(namespace).Create(ctx, og, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create OperatorGroup %s: %w", name, err)
	}

	return nil
}

func createSubscription(ctx context.Context, client dynamic.Interface, namespace, name, channel, source, sourceNamespace string, installPlanApproval string, startingCSV string, labels map[string]interface{}) error {
	subGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "subscriptions",
	}

	spec := map[string]interface{}{
		"channel":             channel,
		"name":                name,
		"source":               source,
		"sourceNamespace":      sourceNamespace,
		"installPlanApproval":  installPlanApproval,
	}

	if startingCSV != "" {
		spec["startingCSV"] = startingCSV
	}

	sub := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "operators.coreos.com/v1alpha1",
			"kind":       "Subscription",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace":  namespace,
				"labels":    labels,
			},
			"spec": spec,
		},
	}

	_, err := client.Resource(subGVR).Namespace(namespace).Create(ctx, sub, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Subscription %s: %w", name, err)
	}

	return nil
}

func waitForInstallPlanRef(ctx context.Context, client dynamic.Interface, namespace, subscriptionName string, timeout time.Duration) (string, error) {
	subGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "subscriptions",
	}

	var installPlanName string
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		sub, err := client.Resource(subGVR).Namespace(namespace).Get(ctx, subscriptionName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		status, found, err := unstructured.NestedMap(sub.Object, "status")
		if err != nil || !found {
			return false, nil
		}

		installPlanRef, found, err := unstructured.NestedMap(status, "installPlanRef")
		if err != nil || !found {
			return false, nil
		}

		name, found, err := unstructured.NestedString(installPlanRef, "name")
		if err != nil || !found || name == "" {
			return false, nil
		}

		installPlanName = name
		return true, nil
	})

	if err != nil {
		return "", fmt.Errorf("timeout waiting for InstallPlan reference: %w", err)
	}

	return installPlanName, nil
}

func approveInstallPlan(ctx context.Context, client dynamic.Interface, namespace, installPlanName string) error {
	ipGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "installplans",
	}

	ip, err := client.Resource(ipGVR).Namespace(namespace).Get(ctx, installPlanName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get InstallPlan %s: %w", installPlanName, err)
	}

	if err := unstructured.SetNestedField(ip.Object, true, "spec", "approved"); err != nil {
		return fmt.Errorf("failed to set approved field: %w", err)
	}

	_, err = client.Resource(ipGVR).Namespace(namespace).Update(ctx, ip, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to approve InstallPlan %s: %w", installPlanName, err)
	}

	return nil
}

func waitForInstalledCSV(ctx context.Context, client dynamic.Interface, namespace, subscriptionName string, timeout time.Duration) (string, error) {
	subGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "subscriptions",
	}

	var csvName string
	err := wait.PollUntilContextTimeout(ctx, 2*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		sub, err := client.Resource(subGVR).Namespace(namespace).Get(ctx, subscriptionName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		status, found, err := unstructured.NestedMap(sub.Object, "status")
		if err != nil || !found {
			return false, nil
		}

		name, found, err := unstructured.NestedString(status, "installedCSV")
		if err != nil || !found || name == "" {
			return false, nil
		}

		csvName = name
		return true, nil
	})

	if err != nil {
		return "", fmt.Errorf("timeout waiting for installedCSV: %w", err)
	}

	return csvName, nil
}

func waitForCSVSucceeded(ctx context.Context, client dynamic.Interface, namespace, csvName string, timeout time.Duration) (string, error) {
	csvGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "clusterserviceversions",
	}

	var phase string
	err := wait.PollUntilContextTimeout(ctx, 5*time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		csv, err := client.Resource(csvGVR).Namespace(namespace).Get(ctx, csvName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		status, found, err := unstructured.NestedMap(csv.Object, "status")
		if err != nil || !found {
			return false, nil
		}

		p, found, err := unstructured.NestedString(status, "phase")
		if err != nil || !found {
			return false, nil
		}

		phase = p
		return phase == "Succeeded", nil
	})

	if err != nil {
		return phase, fmt.Errorf("timeout waiting for CSV %s to reach Succeeded phase (current: %s): %w", csvName, phase, err)
	}

	return phase, nil
}

func getSubscription(ctx context.Context, client dynamic.Interface, namespace, name string) (*unstructured.Unstructured, error) {
	subGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "subscriptions",
	}

	return client.Resource(subGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

func getCSV(ctx context.Context, client dynamic.Interface, namespace, name string) (*unstructured.Unstructured, error) {
	csvGVR := k8sschema.GroupVersionResource{
		Group:    "operators.coreos.com",
		Version:  "v1alpha1",
		Resource: "clusterserviceversions",
	}

	return client.Resource(csvGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}
