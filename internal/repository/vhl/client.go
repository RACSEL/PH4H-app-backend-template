package vhl

import (
	"context"
	"fmt"
	"ips-lacpass-backend/internal/errors"
	"net/http"
)

type VhlRepository struct {
	Client  *http.Client
	BaseURL string
}

func (c *VhlRepository) CreateQr(ctx context.Context, body CreateQrRequest) (*QrData, error) {
	// TODO: To be implemented by the participant
	return nil, &errors.HttpError{
		StatusCode: 500,
		Body:       []map[string]interface{}{{"error": "Not implemented error", "message": "this method is not implemented yet"}},
		Err:        fmt.Errorf("failed to create VHL QR"),
	}
}
