package core

import (
	"context"
	"errors"
	"fmt"
	ipsClient "ips-lacpass-backend/internal/ips/client"
	"ips-lacpass-backend/internal/vhl/client"
	customErrors "ips-lacpass-backend/pkg/errors"
	authMiddleware "ips-lacpass-backend/pkg/middleware"
)

type VhlService struct {
	DefaultClient *client.VhlClient
	Clients       map[string]*client.VhlClient
}

func NewService(r *client.VhlClient) VhlService {
	return VhlService{
		DefaultClient: r,
		Clients:       make(map[string]*client.VhlClient),
	}
}

func (vs *VhlService) getClient(ctx context.Context) *client.VhlClient {
	nodeName := authMiddleware.GetNodeNameFromContext(ctx)
	if nodeName != "" {
		if c, ok := vs.Clients[nodeName]; ok {
			return c
		}
	}
	return vs.DefaultClient
}

func (vs *VhlService) CreateQrCode(ctx context.Context, expiresOn *string, content *string, passCode *string) (*client.QrData, error) {
	if content == nil || *content == "" {
		return nil, &customErrors.HttpError{
			StatusCode: 400,
			Body:       []map[string]interface{}{{"error": "invalid_request", "message": "Content cannot be empty"}},
			Err:        errors.New("content cannot be empty"),
		}
	}

	c := vs.getClient(ctx)
	qrData, err := c.CreateQr(ctx, client.CreateQrRequest{
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
	c := vs.getClient(ctx)
	fmt.Printf("[DEBUG] GetQrIps called with qrData length: %d\n", len(qrData))
	validation, err := c.Validate(ctx, qrData)
	if err != nil {
		fmt.Printf("[ERROR] vs.Client.Validate failed: %v\n", err)
		return nil, err
	}

	if len(validation.ShLinkContent.Url) == 0 {
		fmt.Println("[WARN] Validation returned empty URL")
		return nil, &customErrors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "validation_error", "message": "Validation server did not return a valid access URL"}},
			Err:        errors.New("content cannot be empty"),
		}
	}
	for step := range validation.Status {
		if validation.Status[step].Status != "SUCCESS" {
			fmt.Printf("[WARN] Validation step failed: %s - %s\n", validation.Status[step].Code, validation.Status[step].Description)
			return nil, &customErrors.HttpError{
				StatusCode: 502,
				Body:       []map[string]interface{}{{"error": "unsuccessful_validation", "message": fmt.Sprintf("Validation step unsuccessful. Code: %s. Description: %s", validation.Status[step].Code, validation.Status[step].Description)}},
				Err:        errors.New("Unsuccessful Validation"),
			}
		}
	}

	ipsFetchUrl, err := c.GetIpsUrl(ctx, validation.ShLinkContent.Url, passCode)
	if err != nil {
		fmt.Printf("[ERROR] vs.Client.GetIpsUrl failed: %v\n", err) // <--- Check this!
		return nil, err
	}

	if len(ipsFetchUrl.Files) == 0 || len(ipsFetchUrl.Files[0].Location) == 0 {
		fmt.Println("[WARN] Manifest server returned invalid bundle url")
		return nil, &customErrors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "invalid_manifest_url", "message": "Manifest server returned invalid bundle url."}},
			Err:        errors.New("content cannot be empty"),
		}
	}
	ipsClt := ipsClient.NewClient("", "")
	fmt.Printf("[DEBUG] Fetching IPS Bundle from: %s\n", ipsFetchUrl.Files[0].Location)
	ipsBundle, err := ipsClt.GetIpsBundle(ipsFetchUrl.Files[0].Location)
	if err != nil {
		fmt.Printf("[ERROR] ipsClt.GetIpsBundle failed: %v\n", err)
		return nil, err
	}
	return ipsBundle, nil
}

func (vs *VhlService) GetICVPValidation(ctx context.Context, qrData string) (*client.ICVPQRValidationResponse, error) {
	c := vs.getClient(ctx)
	validationData, err := c.ICVPValidate(ctx, qrData)
	if err != nil {
		return nil, err
	}
	return validationData, nil
}
