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
	Client           *http.Client
	BaseURL          string
	ICVPValidatorUrl string
}

func NewClient(baseURL string, icvpValidatorUrl string) VhlClient {
	return VhlClient{
		Client:           &http.Client{},
		BaseURL:          baseURL,
		ICVPValidatorUrl: icvpValidatorUrl,
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

func (c *VhlClient) ICVPValidate(ctx context.Context, qrData string) (*ICVPQRValidationResponse, error) {
	r := ICVPQrValidationRequest{
		IncludeRaw: true,
		QRData:     qrData,
	}
	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user payload: %w", err)
	}

	vu := fmt.Sprintf("%s/decode/hcert", c.ICVPValidatorUrl)
	req, err := http.NewRequest("POST", vu, bytes.NewBuffer(body))
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: 500,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to create request"}},
			Err:        err,
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)

	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to VHL service"}},
			Err:        err,
		}
	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, &errors.HttpError{
				StatusCode: resp.StatusCode,
				Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to read response body"}},
				Err:        err,
			}
		}
		return nil, &errors.HttpError{
			StatusCode: resp.StatusCode,
			Body:       []map[string]interface{}{{"error": "service_error", "message": string(bodyBytes)}},
			Err:        fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes)),
		}
	}

	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: 500,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to read response body"}},
			Err:        err,
		}
	}
	var valResp ICVPQRValidationResponse
	err = json.Unmarshal(bb, &valResp)
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: 500,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to read parse response body"}},
			Err:        err,
		}
	}
	return &valResp, nil
}
