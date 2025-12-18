package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/redhat/terraform-provider-openshift-operator/internal/client"
	"github.com/redhat/terraform-provider-openshift-operator/internal/resources"
)

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"kubeconfig": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path to kubeconfig file. If not set, uses in-cluster config or KUBECONFIG env var.",
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Kubernetes API server host URL.",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Bearer token for Kubernetes API authentication.",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip TLS certificate verification.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"openshift_operator": resources.ResourceOperator(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"openshift_operator": resources.DataSourceOperator(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var kubeconfig, host, token string
	var insecure bool

	if v, ok := d.GetOk("kubeconfig"); ok {
		kubeconfig = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
	}
	if v, ok := d.GetOk("token"); ok {
		token = v.(string)
	}
	if v, ok := d.GetOk("insecure"); ok {
		insecure = v.(bool)
	}

	cl, err := client.NewClient(kubeconfig, host, token, insecure)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("failed to build Kubernetes client: %w", err))
	}

	return cl, nil
}
