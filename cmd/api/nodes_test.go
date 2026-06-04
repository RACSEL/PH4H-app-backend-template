package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandleGetNodes(t *testing.T) {
	// Create a temporary node-services.json
	configContent := `[
    {
        "id": "node-1",
        "name": "node-1",
        "FHIR_BASE_URL": "http://node-1-fhir",
        "FHIR_MEDIATOR_BASE_URL": "http://node-1-mediator",
        "VHL_BASE_URL": "http://node-1-vhl",
        "ICVP_VALIDATOR_URL": "http://node-1-validator"
    },
    {
        "id": "node-2",
        "name": "node-2",
        "FHIR_BASE_URL": "http://node-2-fhir",
        "FHIR_MEDIATOR_BASE_URL": "http://node-2-mediator",
        "VHL_BASE_URL": "http://node-2-vhl",
        "ICVP_VALIDATOR_URL": "http://node-2-validator"
    }
]`
	err := os.WriteFile("node-services.json", []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove("node-services.json")

	os.Setenv("USE_MULTIPLE_NODES", "1")
	defer os.Unsetenv("USE_MULTIPLE_NODES")

	cfg := LoadConfig()
	app := New(cfg)

	req, err := http.NewRequest("GET", "/nodes", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	app.router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var nodes []map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &nodes)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expected := []map[string]string{
		{"id": "node-1", "name": "node-1"},
		{"id": "node-2", "name": "node-2"},
	}
	if len(nodes) != len(expected) {
		t.Fatalf("Expected %d nodes, got %d", len(expected), len(nodes))
	}

	for i, node := range nodes {
		if node["id"] != expected[i]["id"] || node["name"] != expected[i]["name"] {
			t.Errorf("Expected node %v, got %v", expected[i], node)
		}
	}
}

func TestHandleGetNodesDisabled(t *testing.T) {
	os.Setenv("USE_MULTIPLE_NODES", "0")
	defer os.Unsetenv("USE_MULTIPLE_NODES")

	cfg := LoadConfig()
	app := New(cfg)

	req, err := http.NewRequest("GET", "/nodes", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	app.router.ServeHTTP(rr, req)

	// Since the route is only registered if UseMultipleNodes is true,
	// it should return 404.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
