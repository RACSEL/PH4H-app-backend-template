package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"ips-lacpass-backend/pkg/errors"
	"ips-lacpass-backend/pkg/utils"
	"net/http"
)

type VhlClient struct {
	Client  *http.Client
	BaseURL string
}

func NewClient(baseURL string) VhlClient {
	return VhlClient{
		Client:  &http.Client{},
		BaseURL: baseURL,
	}
}

func (c *VhlClient) CreateQr(ctx context.Context, body CreateQrRequest) (*QrData, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}

func (c *VhlClient) Validate(ctx context.Context, qrData string) (*QRValidationResponse, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}

func (c *VhlClient) GetIpsUrl(ctx context.Context, shLink string, passCode string) (*VhlManifestResponse, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to get document reference"),
	}
}
