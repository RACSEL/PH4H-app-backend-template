package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"ips-lacpass-backend/internal/core"
	"ips-lacpass-backend/internal/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func NewKeycloakClient(baseURL, realm, clientId, clientSecret, emailRedirectUri, emailClientId string, emailLifeSpan int) core.UserRepository {
	return &Client{
		Client:           &http.Client{},
		BaseURL:          baseURL,
		Realm:            realm,
		AdminClientID:    clientId,
		ClientSecret:     clientSecret,
		EmailRedirectURI: emailRedirectUri,
		EmailClientID:    emailClientId,
		EmailLifespan:    emailLifeSpan,
	}
}

func (kc *Client) getAccessToken() (*TokenResponse, error) {
	tokenEndpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.BaseURL, kc.Realm)

	data := url.Values{}
	data.Set("client_id", kc.AdminClientID)
	data.Set("client_secret", kc.ClientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := http.PostForm(tokenEndpoint, data)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Keycloak: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Errorf("failed to close response body: %s", err.Error())
		}
	}(resp.Body)

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

func (kc *Client) sendValidationEmail(ctx context.Context, userID string) error {
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))

	resp, err := kc.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

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

func (kc *Client) CreateUser(ctx context.Context, user core.User, password string) (*core.User, error) {
	r := UserRegistrationRequest{
		Username:  user.Identifier,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Enabled:   true,
		Attributes: map[string][]string{
			"locale":        {user.Locale},
			"document_type": {string(user.DocumentType)},
			"identifier":    {user.Identifier},
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))

	resp, err := kc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body: %v", err)
		}
	}(resp.Body)

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
	fmt.Println("New user ID:", userID)

	go func() {
		newCtx := context.Background()
		if err := kc.sendValidationEmail(newCtx, userID); err != nil {
			// Log the error but don't return it since this is running asynchronously
			fmt.Println("Error sending validation email: %v", err)
		}
	}()
	// Return the first user since we're querying by username/email
	return &core.User{
		ID:           userID,
		Username:     user.Identifier,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Locale:       user.Locale,
		DocumentType: user.DocumentType,
		Identifier:   user.Identifier,
	}, nil
}

func (kc *Client) UpdateUser(ctx context.Context, userUUID string, ur core.UserUpdateRequest) (*core.User, error) {
	// Get UserRepresentation from keycloack
	ku := fmt.Sprintf("%s/admin/realms/%s/users/%s", kc.BaseURL, kc.Realm, userUUID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ku, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	t, err := kc.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))
	resp, err := kc.Client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body: %w", err)
		}
	}(resp.Body)
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
	if ur.FirstName != userRep.FirstName {
		userRep.FirstName = ur.FirstName
		changed = true
	}
	if ur.LastName != userRep.LastName {
		userRep.LastName = ur.LastName
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.AccessToken))

	resp, err = kc.Client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body: %w", err)
		}
	}(resp.Body)
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
	return &core.User{
		ID:           userRep.ID,
		Username:     userRep.Username,
		Email:        userRep.Email,
		FirstName:    ur.FirstName,
		LastName:     ur.LastName,
		Locale:       userRep.Attributes["locale"][0],
		DocumentType: core.AllowedDocumenTypes[userRep.Attributes["document_type"][0]],
		Identifier:   userRep.Attributes["identifier"][0],
	}, nil
}
