package handler

import (
	"encoding/json"
	"errors"
	"ips-lacpass-backend/internal/ips/core"
	errors2 "ips-lacpass-backend/pkg/errors"
	"ips-lacpass-backend/pkg/utils"
	"net/http"
)

type ICVPDataResponse struct {
	Data    string                 `json:"data"`
	Payload map[string]interface{} `json:"payload"`
}

type MergeIPSRequest struct {
	CurrentIPS map[string]interface{} `json:"current_ips"`
	NewIPS     map[string]interface{} `json:"new_ips"`
}

type Handler struct {
	IpsService *core.IpsService
}

func NewHandler(s *core.IpsService) *Handler {
	return &Handler{
		IpsService: s,
	}
}

// GetIPS godoc
//
//	@Summary		Fetch IPS from national node.
//	@Description	Fetch IPS from national node using session access token user identifier.
//	@Tags			IPS FHIR
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Success		200	{object}	any
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/ips [get]
func (ih *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ips, err := ih.IpsService.GetIps(ctx)
	if err != nil {
		var httpErr *errors2.HttpError
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

// MergeIPS godoc
//
//	@Summary		Merge two IPS bundles into a unified IPS.
//	@Description	Merge two FHIR R4 IPS bundles into a single one, removing reduncancy.
//	@Tags			IPS FHIR
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Param			data	body		MergeIPSRequest	true	"IPS bundles to merge"
//
//	@Success		200		{object}	any
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/ips/merge [post]
func (ih *Handler) Merge(w http.ResponseWriter, r *http.Request) {
	var body MergeIPSRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	mi, err := ih.IpsService.MergeIPS(ctx, body.CurrentIPS, body.NewIPS)
	if err != nil {
		// TODO Do correct error handling
		http.Error(w, "Failed to merge IPS", http.StatusInternalServerError)
	}

	res, err := json.Marshal(mi)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// GetICVP godoc
//
//	@Summary		Generate ICVP certificate from an IPS
//	@Description	Generate ICVP vaccination certificate using the id of an IPS and optionally the id of an immunization.
//	@Tags			IPS FHIR
//	@Produce		json
//
//	@Security		ApiKeyAuth
//
//	@Param			bundleId		query		string	true	"IPS bundle id"
//	@Param			immunizationId	query		string	false	"Immunization id"
//
//	@Success		200				{object}	ICVPDataResponse
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/ips/icvp [get]
func (ih *Handler) GetICVP(w http.ResponseWriter, r *http.Request) {
	bundleId := r.URL.Query().Get("bundleId")
	if bundleId == "" {
		http.Error(w, "bundleId query parameter is required", http.StatusBadRequest)
		return
	}
	immunizationId := r.URL.Query().Get("immunizationId")

	ctx := r.Context()
	var immunizationIdPtr *string
	if immunizationId != "" {
		immunizationIdPtr = &immunizationId
	}

	icvp, err := ih.IpsService.GetIpsICVP(ctx, bundleId, immunizationIdPtr)
	if err != nil {
		var httpErr *errors2.HttpError
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

	decodedPayload, err := utils.DecodeHCert(icvp)
	if err != nil {
		http.Error(w, "Failed to decode hcert", http.StatusInternalServerError)
		return
	}

	response := ICVPDataResponse{Data: icvp, Payload: decodedPayload}
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
