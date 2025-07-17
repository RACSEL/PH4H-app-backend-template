package keycloak

import (
	"net/http"
)

type Client struct {
	Client           *http.Client
	BaseURL          string
	Realm            string
	AdminClientID    string
	ClientSecret     string
	EmailRedirectURI string
	EmailLifespan    int
	EmailClientID    string
}

type UserRegistration struct {
	ID          string              `json:"id,omitempty"`
	Username    string              `json:"username"`
	Email       string              `json:"email"`
	FirstName   string              `json:"firstName,omitempty"`
	LastName    string              `json:"lastName,omitempty"`
	Attributes  map[string][]string `json:"attributes,omitempty"`
	Credentials []Credential        `json:"credentials,omitempty"`
}

type Credential struct {
	Type        string `json:"type"`
	Value       string `json:"value,omitempty"`
	Temporary   bool   `json:"temporary"`
	CreatedDate int64  `json:"createdDate,omitempty"`
}

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
}

type UserRegistrationRequest struct {
	Username    string              `json:"username"`
	Email       string              `json:"email"`
	FirstName   string              `json:"firstName"`
	LastName    string              `json:"lastName"`
	Enabled     bool                `json:"enabled"`
	Attributes  map[string][]string `json:"attributes"`
	Credentials []UserCredential    `json:"credentials"`
}

type UserCredential struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}
