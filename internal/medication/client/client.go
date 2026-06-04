package client

import (
	"encoding/json"
	"fmt"
	"io"
	"ips-lacpass-backend/pkg/errors"
	"ips-lacpass-backend/pkg/utils"
	"net/http"
)

type MedicationClient struct {
	Client          *http.Client
	BaseURL         string
	MediatorBaseURL string
}

func NewClient(baseURL string, mediatorBaseURL string) *MedicationClient {
	return &MedicationClient{
		Client:          &http.Client{},
		BaseURL:         baseURL,
		MediatorBaseURL: mediatorBaseURL,
	}
}

func request(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to FHIR service"}},
			Err:        fmt.Errorf("failed to send request: %w", err),
		}
	}

	return resp, nil
}

func (c *MedicationClient) GetDocumentReference(identifier string) (*Bundle, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}

func (c *MedicationClient) GetBundle(uri string) (map[string]interface{}, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}

func (c *MedicationClient) GetMEOW(idBundle string, medicationStatementId *string) (string, error) {
	// TODO: To be implemented by the participant
	return "", &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}
