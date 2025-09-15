package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"ips-lacpass-backend/internal/vhl/core"
	customErrors "ips-lacpass-backend/pkg/errors"
	"ips-lacpass-backend/pkg/utils"
	"net/http"
)

type Handler struct {
	Service *core.VhlService
}

func NewHandler(s *core.VhlService) *Handler {
	return &Handler{
		Service: s,
	}
}

type VhlRequest struct {
	ExpiresOn string `json:"expires_on,omitempty"`
	Content   string `json:"content,required"`
	PassCode  string `json:"pass_code,omitempty"`
}

type VhlGetRequest struct {
	Data     string `json:"data,required"`
	PassCode string `json:"pass_code,omitempty"`
}

type VhlResponse struct {
	Data    string                 `json:"data"`
	Payload map[string]interface{} `json:"payload"`
}

// Create QR data godoc
//
//	@Summary		Create QR data.
//	@Description	Create QR data from VHL issuance.
//	@Tags			IPS FHIR
//	@Accept			json
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Param			data	body		VhlRequest	true	"Data parameters"
//
//	@Success		200		{object}	VhlResponse
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/qr [post]
func (vh *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO check if user is authenticated and has the permission to create a QR code

	// TODO throw correct error body
	var body VhlRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qr, err := vh.Service.CreateQrCode(ctx, &body.ExpiresOn, &body.Content, &body.PassCode)
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

	decodedPayload, err := utils.DecodeHCert(qr.Value)
	if err != nil {
		fmt.Println("Failed to decode hcert: ", err)
	}

	res, err := json.Marshal(&VhlResponse{
		Data:    qr.Value,
		Payload: decodedPayload,
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

// Get IPS Bundle from valid QR data godoc
//
//	@Summary		Get IPS Bundle with valid VHL QR.
//	@Description	Get IPS Bundle using a valid VHL QR.
//	@Tags			IPS FHIR
//	@Accept			json
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Param			data	body		VhlGetRequest	true	"Data parameters"
//
//	@Success		200		{object}	any
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/qr/fetch [post]
func (vh *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO throw correct error body
	var body VhlGetRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ips, err := vh.Service.GetQrIps(ctx, body.Data, body.PassCode)
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

	res, err := json.Marshal(ips)
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
