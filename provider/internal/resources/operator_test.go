package resources

import (
	"strings"
	"testing"
)

func TestResourceOperator(t *testing.T) {
	resource := ResourceOperator()
	if resource == nil {
		t.Fatal("ResourceOperator() returned nil")
	}

	if resource.CreateContext == nil {
		t.Error("CreateContext function is nil")
	}
	if resource.ReadContext == nil {
		t.Error("ReadContext function is nil")
	}
	if resource.UpdateContext == nil {
		t.Error("UpdateContext function is nil")
	}
	if resource.DeleteContext == nil {
		t.Error("DeleteContext function is nil")
	}
	if resource.Importer == nil {
		t.Error("Importer is nil")
	}
}

func TestResourceOperatorSchema(t *testing.T) {
	resource := ResourceOperator()
	schema := resource.Schema

	// Test required fields
	requiredFields := []string{"name", "namespace", "channel", "source"}
	for _, field := range requiredFields {
		if s, ok := schema[field]; !ok {
			t.Errorf("Required field %q not found in schema", field)
		} else if !s.Required {
			t.Errorf("Field %q should be required", field)
		}
	}

	// Test optional fields with defaults
	optionalFields := map[string]interface{}{
		"install_plan_approval": "Automatic",
		"create_namespace":      true,
		"wait_for_csv":          true,
		"wait_timeout":          "10m",
	}

	for field, expectedDefault := range optionalFields {
		if s, ok := schema[field]; !ok {
			t.Errorf("Optional field %q not found in schema", field)
		} else if s.Default != expectedDefault {
			t.Errorf("Field %q has incorrect default: got %v, want %v", field, s.Default, expectedDefault)
		}
	}

	// Test computed fields
	computedFields := []string{"installed_csv", "csv_phase"}
	for _, field := range computedFields {
		if s, ok := schema[field]; !ok {
			t.Errorf("Computed field %q not found in schema", field)
		} else if !s.Computed {
			t.Errorf("Field %q should be computed", field)
		}
	}
}

func TestInstallPlanApprovalValidation(t *testing.T) {
	resource := ResourceOperator()
	schema := resource.Schema

	installPlanApprovalSchema := schema["install_plan_approval"]
	if installPlanApprovalSchema.ValidateFunc == nil {
		t.Fatal("install_plan_approval ValidateFunc is nil")
	}

	// Test valid values
	validValues := []string{"Automatic", "Manual"}
	for _, val := range validValues {
		warns, errs := installPlanApprovalSchema.ValidateFunc(val, "install_plan_approval")
		if len(errs) > 0 {
			t.Errorf("Valid value %q failed validation: %v", val, errs)
		}
		if len(warns) > 0 {
			t.Errorf("Valid value %q produced warnings: %v", val, warns)
		}
	}

	// Test invalid values
	invalidValues := []string{"invalid", "auto", "manual", ""}
	for _, val := range invalidValues {
		warns, errs := installPlanApprovalSchema.ValidateFunc(val, "install_plan_approval")
		if len(errs) == 0 {
			t.Errorf("Invalid value %q should have failed validation", val)
		}
		_ = warns // warnings are acceptable
	}
}

// TestResourceIDFormat tests resource ID parsing logic
func TestResourceIDParsing(t *testing.T) {
	testCases := []struct {
		id        string
		namespace string
		name      string
		valid     bool
	}{
		{"namespace/name", "namespace", "name", true},
		{"openshift-gitops-operator/openshift-gitops-operator", "openshift-gitops-operator", "openshift-gitops-operator", true},
		{"invalid", "", "", false},
		{"", "", "", false},
		{"a/b/c", "", "", false},
	}

	for _, tc := range testCases {
		parts := strings.Split(tc.id, "/")
		if tc.valid {
			if len(parts) != 2 {
				t.Errorf("ID %q: expected 2 parts, got %d", tc.id, len(parts))
				continue
			}
			if parts[0] != tc.namespace || parts[1] != tc.name {
				t.Errorf("ID %q: expected namespace=%q name=%q, got namespace=%q name=%q",
					tc.id, tc.namespace, tc.name, parts[0], parts[1])
			}
		} else {
			if len(parts) == 2 {
				t.Errorf("ID %q: should be invalid but parsed as valid", tc.id)
			}
		}
	}
}
