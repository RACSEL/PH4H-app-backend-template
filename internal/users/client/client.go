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
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ClientInterface interface {
	SendValidationEmail(ctx context.Context, userID string) error
	CreateUser(ctx context.Context, user map[string]interface{}, password string) (*UserID, error)
	UpdateUser(ctx context.Context, userUUID string, ur map[string]interface{}) (*UserRegistration, error)
}

func NewClient(baseURL, realm, clientId, clientSecret, emailRedirectUri, emailClientId string, emailLifeSpan int) UserClient {
	return UserClient{
		Client:           &http.Client{},
		BaseURL:          baseURL,
		Realm:            realm,
		AdminClientID:    clientId,
		ClientSecret:     clientSecret,
		EmailRedirectURI: emailRedirectUri,
		EmailClientID:    emailClientId,
		EmailLifespan:    emailLifeSpan,
		TokenManager:     &TokenManager{},
	}
}

func fetchToken(kc *UserClient) (*TokenResponse, error) {
	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.BaseURL, kc.Realm)

	data := url.Values{}
	data.Set("client_id", kc.AdminClientID)
	data.Set("client_secret", kc.ClientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := http.PostForm(tokenEndpoint, data)
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to Keycloak service"}},
			Err:        err,
		}
	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to obtain token: %s", string(body))
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	return &token, nil
}

func (kc *UserClient) getAccessToken() (string, error) {
	kc.TokenManager.mu.RLock()

	if time.Now().After(kc.TokenManager.tokenExpiresAt.Add(-1 * time.Minute)) {
		kc.TokenManager.mu.RUnlock()
		kc.TokenManager.mu.Lock()
		defer kc.TokenManager.mu.Unlock()

		if time.Now().After(kc.TokenManager.tokenExpiresAt.Add(-1 * time.Minute)) {
			token, err := fetchToken(kc)
			if err != nil {
				return "", fmt.Errorf("failed to get access token: %w", err)
			}
			kc.TokenManager = &TokenManager{
				token:          token.AccessToken,
				tokenExpiresAt: time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
			}
		}
	} else {
		defer kc.TokenManager.mu.RUnlock()
	}

	return kc.TokenManager.token, nil
}

func (kc *UserClient) SendValidationEmail(ctx context.Context, userID string) error {
	actions := []string{"VERIFY_EMAIL"}
	body, err := json.Marshal(actions)
	if err != nil {
		return fmt.Errorf("failed to marshal actions: %w", err)
	}

	ku := fmt.Sprintf("%s/admin/realms/%s/users/%s/execute-actions-email", kc.BaseURL, kc.Realm, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, ku, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("client_id", kc.EmailClientID)
	q.Add("lifespan", strconv.Itoa(kc.EmailLifespan))
	q.Add("redirect_uri", kc.EmailRedirectURI)
	req.URL.RawQuery = q.Encode()

	t, err := kc.getAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))

	resp, err := kc.Client.Do(req)
	if err != nil {
		return &errors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to Auth service"}},
			Err:        err,
		}
	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return &errors.HttpError{
			StatusCode: resp.StatusCode,
			Body:       []map[string]interface{}{errorResponse},
			Err:        err,
		}
	}

	return nil
}

func (kc *UserClient) CreateUser(ctx context.Context, user map[string]interface{}, password string) (*UserID, error) {
	r := UserRegistrationRequest{
		Username:  user["Identifier"].(string),
		Email:     user["Email"].(string),
		FirstName: user["FirstName"].(string),
		LastName:  user["LastName"].(string),
		Enabled:   true,
		Attributes: map[string][]string{
			"locale":        {user["Locale"].(string)},
			"document_type": {user["DocumentType"].(string)},
			"identifier":    {user["Identifier"].(string)},
		},
		Credentials: []UserCredential{
			{
				Type:      "password",
				Value:     password,
				Temporary: false,
			},
		},
	}

	body, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user payload: %w", err)
	}

	ku := fmt.Sprintf("%s/admin/realms/%s/users", kc.BaseURL, kc.Realm)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ku, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	t, err := kc.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))

	resp, err := kc.Client.Do(req)
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to Auth service"}},
			Err:        err,
		}
	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		return nil, &errors.HttpError{
			StatusCode: resp.StatusCode,
			Body:       []map[string]interface{}{errorResponse},
			Err:        err,
		}
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return nil, fmt.Errorf("user created but Location header is missing")
	}

	parts := strings.Split(location, "/")
	userID := parts[len(parts)-1]
	go func() {
		newCtx := context.Background()
		if err := kc.SendValidationEmail(newCtx, userID); err != nil {
			// Log the error but don't return it since this is running asynchronously
			fmt.Errorf("error sending validation email: %v", err)
		}
	}()
	// Return the first user since we're querying by username/email
	return &UserID{ID: userID}, nil
}

func (kc *UserClient) UpdateUser(ctx context.Context, userUUID string, ur map[string]interface{}) (*UserRegistration, error) {
	// Get UserRepresentation from Keycloak
	ku := fmt.Sprintf("%s/admin/realms/%s/users/%s", kc.BaseURL, kc.Realm, userUUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ku, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	t, err := kc.getAccessToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
	resp, err := kc.Client.Do(req)
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to Auth service"}},
			Err:        err,
		}

	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		return nil, &errors.HttpError{
			StatusCode: resp.StatusCode,
			Body:       []map[string]interface{}{errorResponse},
			Err:        err,
		}
	}
	var userRep UserRegistration
	if err := json.NewDecoder(resp.Body).Decode(&userRep); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}

	// Edit User representation
	// For now, only  firstName and lastName can be edited
	changed := false
	if ur["first_name"].(string) != userRep.FirstName {
		userRep.FirstName = ur["first_name"].(string)
		changed = true
	}
	if ur["last_name"].(string) != userRep.LastName {
		userRep.LastName = ur["last_name"].(string)
		changed = true
	}
	if !changed {
		return nil, &errors.HttpError{
			StatusCode: http.StatusBadRequest,
			Body:       []map[string]interface{}{{"error": "update_user_no_change", "message": "update request does not change user attributes"}},
		}
	}

	//save updated user representation
	body, err := json.Marshal(userRep)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated user payload: %w", err)
	}
	ku = fmt.Sprintf("%s/admin/realms/%s/users/%s", kc.BaseURL, kc.Realm, userUUID)
	req, err = http.NewRequestWithContext(ctx, http.MethodPut, ku, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))

	resp, err = kc.Client.Do(req)
	if err != nil {
		return nil, &errors.HttpError{
			StatusCode: 502,
			Body:       []map[string]interface{}{{"error": "service_unavailable", "message": "Failed to connect to Auth service"}},
			Err:        err,
		}
	}
	defer utils.CloseBody(resp.Body)

	if resp.StatusCode != http.StatusNoContent {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		return nil, &errors.HttpError{
			StatusCode: resp.StatusCode,
			Body:       []map[string]interface{}{errorResponse},
			Err:        err,
		}
	}
	return &userRep, nil
}
