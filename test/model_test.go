package test

import (
	"encoding/json"
	"ips-lacpass-backend/internal/ips/core"
	"os"
	"reflect"
	"testing"
)

func TestBundle_Deserialization(t *testing.T) {
	content, _ := os.ReadFile("test_ips.json")

	var bundle core.Bundle
	err := json.Unmarshal(content, &bundle)
	if err != nil {
		t.Fatalf("Failed to unmarshal sample JSON: %v", err)
	}

	remarshaledJSON, err := json.Marshal(bundle)
	if err != nil {
		t.Fatalf("Failed to marshal bundle struct back to JSON: %v", err)
	}

	var originalMap, remarshaledMap map[string]interface{}

	err = json.Unmarshal(content, &originalMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal original JSON to map: %v", err)
	}

	err = json.Unmarshal(remarshaledJSON, &remarshaledMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal remarshaled JSON to map: %v", err)
	}

	if !reflect.DeepEqual(originalMap, remarshaledMap) {
		originalPretty, _ := json.MarshalIndent(originalMap, "", "  ")
		remarshaledPretty, _ := json.MarshalIndent(remarshaledMap, "", "  ")

		t.Errorf("Serialized JSON is not equal to the original.\nOriginal:\n%s\n\nRemarshaled:\n%s", string(originalPretty), string(remarshaledPretty))
	}
}
