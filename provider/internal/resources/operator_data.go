package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/redhat/terraform-provider-openshift-operator/internal/client"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// DataSourceOperator returns a data source for reading operator status
func DataSourceOperator() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOperatorRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the operator (e.g., 'openshift-gitops-operator').",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace where the operator is installed.",
			},
			// Computed outputs
			"channel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current channel for the operator subscription.",
			},
			"source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Catalog source name for the operator subscription.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pinned version (from startingCSV) if version pinning is enabled.",
			},
			"install_plan_approval": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Install plan approval strategy ('Automatic' or 'Manual').",
			},
			"installed_csv": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the installed ClusterServiceVersion (CSV).",
			},
			"csv_phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current phase of the CSV (e.g., 'Succeeded', 'Installing', 'Failed').",
			},
			"csv_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the installed CSV.",
			},
			"subscription_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the subscription (e.g., 'AtLatestKnown', 'UpgradePending').",
			},
			"current_csv": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current CSV that the subscription is tracking.",
			},
			"installed_csv_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version extracted from the installed CSV name.",
			},
		},
	}
}

func dataSourceOperatorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	dynamicClient := c.Dynamic

	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)

	// Validate inputs
	if namespace == "" {
		return diag.FromErr(fmt.Errorf("namespace cannot be empty"))
	}
	if name == "" {
		return diag.FromErr(fmt.Errorf("name cannot be empty"))
	}

	// Read Subscription
	sub, err := getSubscription(ctx, dynamicClient, namespace, name)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Operator not found",
				Detail:   fmt.Sprintf("Subscription '%s' not found in namespace '%s'. Error: %v", name, namespace, err),
			},
		}
	}

	// Read Subscription spec
	if spec, found, _ := unstructured.NestedMap(sub.Object, "spec"); found {
		if channel, found, _ := unstructured.NestedString(spec, "channel"); found {
			d.Set("channel", channel)
		}
		if source, found, _ := unstructured.NestedString(spec, "source"); found {
			d.Set("source", source)
		}
		if installPlanApproval, found, _ := unstructured.NestedString(spec, "installPlanApproval"); found {
			d.Set("install_plan_approval", installPlanApproval)
		}
		if startingCSV, found, _ := unstructured.NestedString(spec, "startingCSV"); found && startingCSV != "" {
			d.Set("version", startingCSV)
			// Extract version from CSV name format: {name}.v{version}
			if strings.HasPrefix(startingCSV, name+".v") {
				version := strings.TrimPrefix(startingCSV, name+".v")
				d.Set("version", version)
			}
		}
	}

	// Read Subscription status
	if status, found, _ := unstructured.NestedMap(sub.Object, "status"); found {
		if installedCSV, found, _ := unstructured.NestedString(status, "installedCSV"); found && installedCSV != "" {
			d.Set("installed_csv", installedCSV)
			d.Set("current_csv", installedCSV)

			// Extract version from CSV name
			if strings.HasPrefix(installedCSV, name+".v") {
				version := strings.TrimPrefix(installedCSV, name+".v")
				d.Set("installed_csv_version", version)
			}

			// Read CSV details
			csv, err := getCSV(ctx, dynamicClient, namespace, installedCSV)
			if err == nil {
				// Read CSV spec for version
				if csvSpec, found, _ := unstructured.NestedMap(csv.Object, "spec"); found {
					if csvVersion, found, _ := unstructured.NestedString(csvSpec, "version"); found {
						d.Set("csv_version", csvVersion)
					}
				}

				// Read CSV status for phase
				if csvStatus, found, _ := unstructured.NestedMap(csv.Object, "status"); found {
					if phase, found, _ := unstructured.NestedString(csvStatus, "phase"); found {
						d.Set("csv_phase", phase)
					}
				}
			}
		}

		if currentCSV, found, _ := unstructured.NestedString(status, "currentCSV"); found && currentCSV != "" {
			d.Set("current_csv", currentCSV)
		}

		if state, found, _ := unstructured.NestedString(status, "state"); found {
			d.Set("subscription_state", state)
		}
	}

	// Set ID and required fields
	d.SetId(fmt.Sprintf("%s/%s", namespace, name))
	d.Set("namespace", namespace)
	d.Set("name", name)

	return nil
}
