package client

type CredentialType string

const (
	VerifiableHealthLink CredentialType = "VerifiableHealthLink"
	ICVP                 CredentialType = "ICVP"
)

// GenerateWalletLinkRequest represents the request body for generating a wallet link.
type GenerateWalletLinkRequest struct {
	Claims         map[string]interface{} `json:"claims"`
	CredentialType CredentialType         `json:"credentialType"`
	PinRequired    bool                   `json:"pinRequired"`
}

// GenerateWalletLinkResponse represents the response from the wallet service.
type GenerateWalletLinkResponse struct {
	PreAuthorizedCode string `json:"preAuthorizedCode"`
	QrURL             string `json:"qrUrl"`
	CoURL             string `json:"coUrl"`
	Location          string `json:"location"`
}
