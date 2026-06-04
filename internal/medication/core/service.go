package core

import (
	"context"
	"fmt"
	"ips-lacpass-backend/internal/medication/client"
	customErrors "ips-lacpass-backend/pkg/errors"
	authMiddleware "ips-lacpass-backend/pkg/middleware"
	"log/slog"
	"sort"
)

type ClientAdapter interface {
	GetDocumentReference(identifier string) (*client.Bundle, error)
	GetBundle(url string) (map[string]interface{}, error)
	GetMEOW(idBundle string, medicationStatementId *string) (string, error)
}

type Service struct {
	Client ClientAdapter
}

func NewService(client ClientAdapter) *Service {
	return &Service{
		Client: client,
	}
}

func (s *Service) GetMedication(ctx context.Context) (map[string]interface{}, error) {
	userId, err := authMiddleware.GetUserDocIDFromContext(ctx)
	if err != nil {
		return nil, &customErrors.HttpError{
			StatusCode: 401,
			Body:       []map[string]interface{}{{"error": "user_identifier_not_found", "message": "User identifier not found in request context"}},
			Err:        err,
		}
	}

	docRef, err := s.Client.GetDocumentReference(userId)
	if err != nil {
		fmt.Printf("Error fetching document reference: %v\n", err)
		return nil, err
	}
	entries := docRef.Entry
	if len(entries) == 0 {
		return nil, &customErrors.HttpError{
			StatusCode: 404,
			Body:       []map[string]interface{}{{"error": "not_found", "message": "No Medication found for the user"}},
			Err:        fmt.Errorf("no medication found for the user"),
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Resource.Meta.LastUpdated > entries[j].Resource.Meta.LastUpdated
	})

	bundle, err := s.Client.GetBundle(entries[0].Resource.Content[0].Attachment.URL)
	if err != nil {
		return nil, err
	}

	bundleType, ok := bundle["type"]
	if !ok || bundleType == nil {
		slog.Warn("Bundle has no type", "bundle", bundle)
		return bundle, nil
	}
	if bundleType != "document" {
		slog.Warn("Medication is not of required type, expected `document`", "bundle_type", bundleType)
	}
	return bundle, nil
}

func (s *Service) GetMEOW(idBundle string, medicationStatementId *string) (string, error) {
	result, err := s.Client.GetMEOW(idBundle, medicationStatementId)
	if err != nil {
		fmt.Printf("[medications/meow error] client failed bundleId=%s hasMedicationStatementId=%t error=%v\n", idBundle, medicationStatementId != nil, err)
		return "", err
	}

	return result, nil
}
