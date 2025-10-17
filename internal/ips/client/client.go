package client

import (
	"encoding/json"
	"fmt"
	"io"
	"ips-lacpass-backend/pkg/errors"
	"ips-lacpass-backend/pkg/utils"
	"net/http"
)

type ClientInterface interface {
	GetDocumentReference(identifier string) (*Bundle, error)
	GetIpsBundle(url string) (map[string]interface{}, error)
	GetIpsICVP(idBundle string, immunizationId *string) (string, error)
}

type IpsClient struct {
	Client          *http.Client
	BaseURL         string
	MediatorBaseURL string
}

func NewClient(baseURL string, mediatorBaseURL string) IpsClient {
	return IpsClient{
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

	if resp.StatusCode != http.StatusOK {
		return nil, &errors.HttpError{
			StatusCode: resp.StatusCode,
			Body:       []map[string]interface{}{{"error": "fhir_error", "message": resp.Body}},
		}
	}
	return resp, nil
}

func (c *IpsClient) GetDocumentReference(identifier string) (*Bundle, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}

func (c *IpsClient) GetIpsBundle(url string) (map[string]interface{}, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}

func (c *IpsClient) GetIpsICVP(idBundle string, immunizationId *string) (string, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}