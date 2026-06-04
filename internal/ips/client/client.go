package client

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"ips-lacpass-backend/pkg/errors"
	"ips-lacpass-backend/pkg/utils"
	"log/slog"
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
	slog.Info("FHIR Request", "method", req.Method, "url", req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("FHIR Request failed", "method", req.Method, "url", req.URL.String(), "error", err)
		return nil, &errors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to FHIR service"}},
			Err:        fmt.Errorf("failed to send request: %w", err),
		}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("FHIR Request error response", "method", req.Method, "url", req.URL.String(), "status", resp.StatusCode, "body", string(body))
		return nil, &errors.HttpError{
			StatusCode: resp.StatusCode,
			Body:       []map[string]interface{}{fhirErrorBody(body)},
			Err:        fmt.Errorf("FHIR service returned status %d", resp.StatusCode),
		}
	}
	slog.Info("FHIR Request success", "method", req.Method, "url", req.URL.String(), "status", resp.StatusCode)
	return resp, nil
}

func fhirErrorBody(body []byte) map[string]interface{} {
	errBody := map[string]interface{}{
		"error": "fhir_error",
	}

	var upstreamBody interface{}
	if err := json.Unmarshal(body, &upstreamBody); err == nil {
		errBody["message"] = "FHIR service returned an error"
		errBody["fhir_error"] = upstreamBody
		return errBody
	}

	errBody["message"] = string(body)
	return errBody
}

func preserveFhirError(err error) error {
	var httpErr *errors.HttpError
	if stderrors.As(err, &httpErr) {
		return err
	}
	return &errors.HttpError{
		StatusCode: 502,
		Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to FHIR service"}},
		Err:        err,
	}
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
	return "", &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}
