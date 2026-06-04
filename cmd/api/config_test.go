package main

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary node-services.json in the current directory for the test
	configContent := `[
    {
        "id": "test-node",
        "name": "test-node",
        "FHIR_BASE_URL": "http://test-fhir",
        "FHIR_MEDIATOR_BASE_URL": "http://test-mediator",
        "VHL_BASE_URL": "http://test-vhl",
        "ICVP_VALIDATOR_URL": "http://test-validator"
    }
]`
	err := os.WriteFile("node-services.json", []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove("node-services.json")

	// Set environment variables
	os.Setenv("USE_MULTIPLE_NODES", "1")
	defer os.Unsetenv("USE_MULTIPLE_NODES")

	cfg := LoadConfig()

	if !cfg.UseMultipleNodes {
		t.Errorf("Expected UseMultipleNodes to be true")
	}

	if len(cfg.Nodes) == 0 {
		t.Fatalf("Expected at least one node in config")
	}

	if cfg.Nodes[0].Name != "test-node" {
		t.Errorf("Expected test-node, got %s", cfg.Nodes[0].Name)
	}

	if cfg.Nodes[0].ID != "test-node" {
		t.Errorf("Expected test-node, got %s", cfg.Nodes[0].ID)
	}

	if cfg.Nodes[0].FhirBaseUrl != "http://test-fhir" {
		t.Errorf("Expected http://test-fhir, got %s", cfg.Nodes[0].FhirBaseUrl)
	}
}
