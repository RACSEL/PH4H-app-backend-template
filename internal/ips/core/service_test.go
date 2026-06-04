package core

import (
	"context"
	"ips-lacpass-backend/internal/ips/client"
	"ips-lacpass-backend/pkg/middleware"
	"testing"
)

func TestNodeSelection(t *testing.T) {
	defaultRepo := client.NewClient("http://default-fhir", "http://default-mediator")
	node1Repo := client.NewClient("http://node1-fhir", "http://node1-mediator")
	
	service := NewService(&defaultRepo)
	service.Repositories["node1"] = &node1Repo

	t.Run("Default node", func(t *testing.T) {
		ctx := context.Background()
		repo := service.getClient(ctx)
		if repo.BaseURL != "http://default-fhir" {
			t.Errorf("Expected default-fhir, got %s", repo.BaseURL)
		}
	})

	t.Run("Node1 selection", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), middleware.NodeNameKey, "node1")
		repo := service.getClient(ctx)
		if repo.BaseURL != "http://node1-fhir" {
			t.Errorf("Expected node1-fhir, got %s", repo.BaseURL)
		}
	})

	t.Run("Unknown node fallback", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), middleware.NodeNameKey, "unknown")
		repo := service.getClient(ctx)
		if repo.BaseURL != "http://default-fhir" {
			t.Errorf("Expected default-fhir, got %s", repo.BaseURL)
		}
	})
}
