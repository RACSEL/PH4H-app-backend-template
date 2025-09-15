package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"ips-lacpass-backend/pkg/errors"
	"ips-lacpass-backend/pkg/utils"
	"net/http"
)

type ClientInterface interface {
	GenerateWalletLink(ctx context.Context, claims map[string]interface{}) (*GenerateWalletLinkResponse, error)
}

type WalletClient struct {
	Client     *http.Client
	BaseURL    string
	Identifier string
	APIKey     string
}

func NewClient(baseURL string, identifier string, apiKey string) WalletClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return WalletClient{
		Client:     &http.Client{Transport: tr},
		BaseURL:    baseURL,
		Identifier: identifier,
		APIKey:     apiKey,
	}
}

func (c *WalletClient) GenerateWalletLink(ctx context.Context, claims map[string]interface{}, credentialType CredentialType) (*GenerateWalletLinkResponse, error) {
	url := fmt.Sprintf("%s/credentials/%s", c.BaseURL, c.Identifier)

	reqBody := GenerateWalletLinkRequest{
		Claims:         claims,
		CredentialType: credentialType,
		PinRequired:    false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: http.StatusInternalServerError,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to marshal request body"}},
			Err:        err,
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: http.StatusInternalServerError,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to create request"}},
			Err:        err,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: http.StatusServiceUnavailable,
			Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to wallet service"}},
			Err:        err,
		}
	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, &errors.HttpError{
			StatusCode: resp.StatusCode,
			Body:       []map[string]interface{}{{"error": "wallet_service_error", "message": "Wallet service returned an error"}},
			Err:        fmt.Errorf("wallet service returned status code %d", resp.StatusCode),
		}
	}

	var walletResponse GenerateWalletLinkResponse
	if err := json.NewDecoder(resp.Body).Decode(&walletResponse); err != nil {
		return nil, &errors.HttpError{
			StatusCode: http.StatusInternalServerError,
			Body:       []map[string]interface{}{{"error": "internal_error", "message": "Failed to decode response"}},
			Err:        fmt.Errorf("failed to decode response: %w", err),
		}
	}

	return &walletResponse, nil
}
