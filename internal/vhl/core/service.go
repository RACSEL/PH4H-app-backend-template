package core

import (
	"context"
	"errors"
	"fmt"
	ipsClient "ips-lacpass-backend/internal/ips/client"
	"ips-lacpass-backend/internal/vhl/client"
	customErrors "ips-lacpass-backend/pkg/errors"
)

type VhlService struct {
	Client *client.VhlClient
}

func NewService(r *client.VhlClient) VhlService {
	return VhlService{
		Client: r,
	}
}

func (vs *VhlService) CreateQrCode(ctx context.Context, expiresOn *string, content *string, passCode *string) (*client.QrData, error) {
	if content == nil || *content == "" {
		return nil, &customErrors.HttpError{
			StatusCode: 400,
			Body:       []map[string]interface{}{{"error": "invalid_request", "message": "Content cannot be empty"}},
			Err:        errors.New("content cannot be empty"),
		}
	}

	qrData, err := vs.Client.CreateQr(ctx, client.CreateQrRequest{
		JsonContent: *content,
		ExpiresOn:   *expiresOn,
		PassCode:    *passCode,
	})
	if err != nil {
		return nil, err
	}
	return qrData, nil
}

func (vs *VhlService) GetQrIps(ctx context.Context, qrData string, passCode string) (map[string]any, error) {
	validation, err := vs.Client.Validate(ctx, qrData)
	if err != nil {
		return nil, err
	}

	if len(validation.ShLinkContent.Url) == 0 {
		return nil, &customErrors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "validation_error", "message": "Validation server did not return a valid access URL"}},
			Err:        errors.New("content cannot be empty"),
		}
	}
	for step := range validation.Status {
		if validation.Status[step].Status != "SUCCESS" {
			return nil, &customErrors.HttpError{
				StatusCode: 502,
				Body:       []map[string]interface{}{{"error": "unsuccessful_validation", "message": fmt.Sprintf("Validation step unsuccessful. Code: %s. Description: %s", validation.Status[step].Code, validation.Status[step].Description)}},
				Err:        errors.New("Unsuccessful Validation"),
			}
		}
	}

	ipsFetchUrl, err := vs.Client.GetIpsUrl(ctx, validation.ShLinkContent.Url, passCode)

	if len(ipsFetchUrl.Files) == 0 || len(ipsFetchUrl.Files[0].Location) == 0 {
		return nil, &customErrors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "invalid_manifest_url", "message": "Manifest server returned invalid bundle url."}},
			Err:        errors.New("content cannot be empty"),
		}
	}
	ipsClt := ipsClient.NewClient("")
	ipsBundle, err := ipsClt.GetIpsBundle(ipsFetchUrl.Files[0].Location)
	if err != nil {
		return nil, err
	}
	return ipsBundle, nil
}
