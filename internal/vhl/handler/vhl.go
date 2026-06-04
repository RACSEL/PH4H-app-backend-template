package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"ips-lacpass-backend/internal/vhl/core"
	customErrors "ips-lacpass-backend/pkg/errors"
	"ips-lacpass-backend/pkg/utils"
	"log"
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

type ICVPValidateRequest struct {
	Data string `json:"data,require"`
}

type VhlResponse struct {
	Data    string                 `json:"data"`
	Payload map[string]interface{} `json:"payload"`
}

// Create Create QR data from VHL issuance.
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

// Get IPS Bundle using a valid VHL QR.
func (vh *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO throw correct error body
	var body VhlGetRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("JSON Decode Error: %v", err)
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
			log.Printf("Service Layer Error: %v", err)
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

// Validate ICVP data for ICVP that dont come from a IPS
func (vh *Handler) Validate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO throw correct error body
	var body ICVPValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("[DEBUG] /qr/validate request received data_len=%d\n", len(body.Data))

	icvpValidationResponseData, err := vh.Service.GetICVPValidation(ctx, body.Data)
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

	res, err := json.Marshal(icvpValidationResponseData)
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
