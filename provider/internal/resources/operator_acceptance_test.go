//go:build acceptance

package resources

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/redhat/terraform-provider-openshift-operator/internal/client"
	"github.com/redhat/terraform-provider-openshift-operator/internal/provider"
)

// Acceptance tests require TF_ACC=1 environment variable and valid cluster access
func TestAccOperator_Basic(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Skipping acceptance test. Set TF_ACC=1 to enable.")
	}

	resourceName := "openshift_operator.test"
	testNamespace := "terraform-test-operator"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOperatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOperatorConfigBasic(testNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOperatorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "namespace", testNamespace),
					resource.TestCheckResourceAttr(resourceName, "channel", "latest"),
					resource.TestCheckResourceAttrSet(resourceName, "installed_csv"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	// Verify that we have required environment variables
	if os.Getenv("KUBECONFIG") == "" && os.Getenv("KUBE_HOST") == "" {
		t.Fatal("KUBECONFIG or KUBE_HOST environment variable must be set for acceptance tests")
	}
}

func testAccCheckOperatorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		return nil
	}
}

func testAccCheckOperatorDestroy(s *terraform.State) error {
	// Verify that the operator resources have been deleted
	// This is a placeholder - actual implementation would check Kubernetes resources
	return nil
}

func testAccOperatorConfigBasic(namespace string) string {
	return fmt.Sprintf(`
provider "openshift" {
  kubeconfig = "%s"
}

resource "openshift_operator" "test" {
  name      = "openshift-gitops-operator"
  namespace = "%s"
  channel   = "latest"
  source    = "redhat-operators"

  wait_for_csv = false  # Speed up tests
}
`, os.Getenv("KUBECONFIG"), namespace)
}

var testAccProviders map[string]*schema.Provider

func init() {
	testAccProviders = map[string]*schema.Provider{
		"openshift": provider.New(),
	}
}
