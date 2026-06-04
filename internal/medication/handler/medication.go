package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	errors2 "ips-lacpass-backend/pkg/errors"
	"net/http"
)

type ServiceAdapter interface {
	GetMedication(ctx context.Context) (map[string]interface{}, error)
	GetMEOW(idBundle string, medicationStatementId *string) (string, error)
}
type Handler struct {
	Service ServiceAdapter
}

func NewHandler(service ServiceAdapter) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ips, err := h.Service.GetMedication(ctx)
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

func (h *Handler) GetMeow(w http.ResponseWriter, r *http.Request) {
	bundleId := r.URL.Query().Get("bundleId")
	if bundleId == "" {
		fmt.Printf("[medications/meow error] missing bundleId path=%s remoteAddr=%s\n", r.URL.Path, r.RemoteAddr)
		http.Error(w, "bundleId query parameter is required", http.StatusBadRequest)
		return
	}

	medicationStatementId := r.URL.Query().Get("medicationStatementId")
	var medicationStatementIdPtr *string
	if medicationStatementId != "" {
		medicationStatementIdPtr = &medicationStatementId
	}

	ips, err := h.Service.GetMEOW(bundleId, medicationStatementIdPtr)
	if err != nil {
		fmt.Printf("[medications/meow error] service failed bundleId=%s hasMedicationStatementId=%t error=%v\n", bundleId, medicationStatementIdPtr != nil, err)
		var httpErr *errors2.HttpError
		if errors.As(err, &httpErr) {
			res, err := json.Marshal(httpErr.Body)
			if err != nil {
				fmt.Printf("[medications/meow error] failed to encode error response bundleId=%s error=%v\n", bundleId, err)
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(httpErr.StatusCode)
			_, err = w.Write(res)
			if err != nil {
				fmt.Printf("[medications/meow error] failed to write error response bundleId=%s statusCode=%d error=%v\n", bundleId, httpErr.StatusCode, err)
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		}
		return
	}

	res, err := json.Marshal(map[string]string{
		"data": ips,
	})
	if err != nil {
		fmt.Printf("[medications/meow error] failed to encode success response bundleId=%s error=%v\n", bundleId, err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		fmt.Printf("[medications/meow error] failed to write success response bundleId=%s error=%v\n", bundleId, err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
