package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNodeNameFromHeader(t *testing.T) {
	var got string
	handler := NodeNameFromHeader(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = GetNodeNameFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/ips", nil)
	req.Header.Set(NodeNameHeader, "node-1")

	handler.ServeHTTP(httptest.NewRecorder(), req)

	if got != "node-1" {
		t.Fatalf("expected node-1, got %q", got)
	}
}

func TestNodeNameFromHeaderIgnoresQueryParam(t *testing.T) {
	var got string
	handler := NodeNameFromHeader(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = GetNodeNameFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/ips?node_name=node-1", nil)

	handler.ServeHTTP(httptest.NewRecorder(), req)

	if got != "" {
		t.Fatalf("expected empty node name, got %q", got)
	}
}
