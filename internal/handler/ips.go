package handler

import (
	"encoding/json"
	"errors"
	"ips-lacpass-backend/internal/core"
	errors2 "ips-lacpass-backend/internal/errors"
	"net/http"
)

type IpsHandler struct {
	FhirService core.FhirService
}

func NewIpsHandler(s core.FhirService) *IpsHandler {
	return &IpsHandler{
		FhirService: s,
	}
}

// GetIPS godoc
//
//	@Summary	    Fetch IPS from national node.
//	@Description	Fetch IPS from national node using session access token user identifier.
//	@Tags			IPS FHIR
//	@Produce		json
//
// @Security ApiKeyAuth
//
//	@Success		200		{object}	fhir.Bundle
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Router			/ips [get]
func (ih *IpsHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ips, err := ih.FhirService.GetIps(ctx)
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
