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
	spec, found, err := unstructured.NestedMap(sub.Object, "spec")
	if err == nil && found {
		channel, found, _ := unstructured.NestedString(spec, "channel")
		if found {
			_ = d.Set("channel", channel)
		}
		source, found, _ := unstructured.NestedString(spec, "source")
		if found {
			_ = d.Set("source", source)
		}
		installPlanApproval, found, _ := unstructured.NestedString(spec, "installPlanApproval")
		if found {
			_ = d.Set("install_plan_approval", installPlanApproval)
		}
		startingCSV, found, _ := unstructured.NestedString(spec, "startingCSV")
		if found && startingCSV != "" {
			_ = d.Set("version", startingCSV)
			// Extract version from CSV name format: {name}.v{version}
			if strings.HasPrefix(startingCSV, name+".v") {
				version := strings.TrimPrefix(startingCSV, name+".v")
				_ = d.Set("version", version)
			}
		}
	}

	// Read Subscription status
	status, found, err := unstructured.NestedMap(sub.Object, "status")
	if err == nil && found {
		installedCSV, found, _ := unstructured.NestedString(status, "installedCSV")
		if found && installedCSV != "" {
			_ = d.Set("installed_csv", installedCSV)
			_ = d.Set("current_csv", installedCSV)

			// Extract version from CSV name
			if strings.HasPrefix(installedCSV, name+".v") {
				version := strings.TrimPrefix(installedCSV, name+".v")
				_ = d.Set("installed_csv_version", version)
			}

			// Read CSV details
			csv, err := getCSV(ctx, dynamicClient, namespace, installedCSV)
			if err == nil {
				// Read CSV spec for version
				csvSpec, found, _ := unstructured.NestedMap(csv.Object, "spec")
				if found {
					csvVersion, found, _ := unstructured.NestedString(csvSpec, "version")
					if found {
						_ = d.Set("csv_version", csvVersion)
					}
				}

				// Read CSV status for phase
				csvStatus, found, _ := unstructured.NestedMap(csv.Object, "status")
				if found {
					phase, found, _ := unstructured.NestedString(csvStatus, "phase")
					if found {
						_ = d.Set("csv_phase", phase)
					}
				}
			}
		}

		currentCSV, found, _ := unstructured.NestedString(status, "currentCSV")
		if found && currentCSV != "" {
			_ = d.Set("current_csv", currentCSV)
		}

		state, found, _ := unstructured.NestedString(status, "state")
		if found {
			_ = d.Set("subscription_state", state)
		}
	}

	// Set ID and required fields
	d.SetId(fmt.Sprintf("%s/%s", namespace, name))
	_ = d.Set("namespace", namespace)
	_ = d.Set("name", name)

	return nil
}
