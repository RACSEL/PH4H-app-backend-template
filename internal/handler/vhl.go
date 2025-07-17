package handler

import (
	"encoding/json"
	"errors"
	"ips-lacpass-backend/internal/core"
	customErrors "ips-lacpass-backend/internal/errors"
	"net/http"
)

type VhlHandler struct {
	VhlService core.VhlService
}

func NewVhlHandler(s core.VhlService) *VhlHandler {
	return &VhlHandler{
		VhlService: s,
	}
}

type VhlRequest struct {
	ExpiresOn string `json:"expires_on,omitempty"`
	Content   string `json:"content,required"`
	PassCode  string `json:"pass_code,omitempty"`
}

type VhlResponse struct {
	Data string `json:"data"`
}

// Create QR data godoc
//
//	@Summary	    Create QR data.
//	@Description	Create QR data from VHL issuance.
//	@Tags			IPS FHIR
//	@Accept			json
//	@Produce		json
//
//	@Param			data	body		VhlRequest	true	"Data parameters"
//
//	@Security ApiKeyAuth
//
//	@Success		200		{object}	VhlResponse
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/qr [post]
func (vh *VhlHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO check if user is authenticated and has the permission to create a QR code

	// TODO throw correct error body
	var body VhlRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qr, err := vh.VhlService.CreateQrCode(ctx, &body.ExpiresOn, &body.Content, &body.PassCode)
	if err != nil {
		var httpErr *customErrors.HttpError
		if errors.As(err, &httpErr) {
			res, err := json.Marshal(httpErr.Body)
			if err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(httpErr.StatusCode)
			_, err = w.Write(res)
			if err != nil {
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		return
	}

	res, err := json.Marshal(&VhlResponse{
		Data: qr.Value,
	})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
