package resources

import (
	"testing"
)

func TestDataSourceOperator(t *testing.T) {
	dataSource := DataSourceOperator()
	if dataSource == nil {
		t.Fatal("DataSourceOperator() returned nil")
	}

	if dataSource.ReadContext == nil {
		t.Error("ReadContext function is nil")
	}
}

func TestDataSourceOperatorSchema(t *testing.T) {
	dataSource := DataSourceOperator()
	schema := dataSource.Schema

	// Test required fields
	requiredFields := []string{"name", "namespace"}
	for _, field := range requiredFields {
		if s, ok := schema[field]; !ok {
			t.Errorf("Required field %q not found in schema", field)
		} else if !s.Required {
			t.Errorf("Field %q should be required", field)
		}
	}

	// Test computed fields
	computedFields := []string{
		"channel",
		"source",
		"version",
		"install_plan_approval",
		"installed_csv",
		"csv_phase",
		"csv_version",
		"subscription_state",
		"current_csv",
		"installed_csv_version",
	}

	for _, field := range computedFields {
		if s, ok := schema[field]; !ok {
			t.Errorf("Computed field %q not found in schema", field)
		} else if !s.Computed {
			t.Errorf("Field %q should be computed", field)
		}
	}
}
