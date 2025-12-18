package resources

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOperator() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOperatorCreate,
		ReadContext:   resourceOperatorRead,
		UpdateContext: resourceOperatorUpdate,
		DeleteContext: resourceOperatorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the operator (e.g., 'openshift-gitops-operator').",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace where the operator will be installed.",
			},
			"channel": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Channel for the operator subscription (e.g., 'latest', 'stable').",
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Catalog source name (e.g., 'redhat-operators', 'certified-operators').",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specific version to install (e.g., '1.18.2'). If set, installPlanApproval is automatically set to 'Manual'.",
			},
			"install_plan_approval": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Automatic",
				Description: "Install plan approval strategy ('Automatic' or 'Manual').",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Automatic" && v != "Manual" {
						errs = append(errs, fmt.Errorf("%q must be either 'Automatic' or 'Manual', got: %s", key, v))
					}
					return
				},
			},
			"create_namespace": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to create the namespace if it doesn't exist.",
			},
			"namespace_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Labels to apply to the namespace.",
			},
			"operator_group_target_namespaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Target namespaces for the OperatorGroup. Empty list means cluster-wide.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"wait_for_csv": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to wait for the CSV to be in 'Succeeded' phase.",
			},
			"wait_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "10m",
				Description: "Timeout for waiting for CSV to succeed (e.g., '10m', '1h').",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Labels to apply to operator resources.",
			},
			// Computed outputs
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
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

// CRUD operations are implemented in operator_impl.go
