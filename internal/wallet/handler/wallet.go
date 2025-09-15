package handler

import (
	"encoding/json"
	"errors"
	"ips-lacpass-backend/internal/wallet/client"
	"ips-lacpass-backend/internal/wallet/core"
	customErrors "ips-lacpass-backend/pkg/errors"
	"net/http"
)

type Handler struct {
	WalletService *core.WalletService
}

func NewHandler(s *core.WalletService) *Handler {
	return &Handler{
		WalletService: s,
	}
}

type GenerateWalletLinkRequest struct {
	Claims         map[string]interface{} `json:"claims"`
	CredentialType client.CredentialType  `json:"credentialType"`
}

// GenerateWalletLink godoc
//
//	@Summary		Generate a wallet link.
//	@Description	Generate a new wallet link with the given claims. Must enable wallet in config.
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Param			data	body		GenerateWalletLinkRequest	true	"Claims for the wallet link"
//
//	@Success		200		{object}	client.GenerateWalletLinkResponse
//	@Failure		400
//	@Failure		500
//	@Router			/wallet/generate-link [post]
func (h *Handler) GenerateWalletLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var reqBody GenerateWalletLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error": "invalid_request_body"}`, http.StatusBadRequest)
		return
	}

	if reqBody.CredentialType != client.VerifiableHealthLink && reqBody.CredentialType != client.ICVP {
		http.Error(w, `{"error": "invalid_credential_type", "message": "credentialType must be either VerifiableHealthLink or ICVP"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	walletResponse, err := h.WalletService.GenerateWalletLink(ctx, reqBody.Claims, reqBody.CredentialType)
	if err != nil {
		var httpErr *customErrors.HttpError
		if errors.As(err, &httpErr) {
			res, err := json.Marshal(httpErr.Body)
			if err != nil {
				http.Error(w, `{"error": "internal_server_error"}`, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(httpErr.StatusCode)
			_, err = w.Write(res)
			if err != nil {
				http.Error(w, `{"error": "internal_server_error"}`, http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, `{"error": "internal_server_error"}`, http.StatusInternalServerError)
		}
		return
	}

	res, err := json.Marshal(walletResponse)
	if err != nil {
		http.Error(w, `{"error": "internal_server_error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, `{"error": "internal_server_error"}`, http.StatusInternalServerError)
		return
	}
}
