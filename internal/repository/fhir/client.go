package fhir

import (
	"fmt"
	"ips-lacpass-backend/internal/errors"
	"net/http"
)

type FhirRepository struct {
	Client  *http.Client
	BaseURL string
}

func (c *FhirRepository) GetDocumentReference(identifier string) (*Bundle, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}

func (c *FhirRepository) GetIpsBundle(url string) (map[string]interface{}, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get Ips Bundle"),
	}
}
