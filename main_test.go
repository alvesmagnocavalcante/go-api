package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name, method, path, key, value string
		status                         int
	}{
		{"root", "GET", "/", "", "", http.StatusOK},
		{"health", "GET", "/health", "status", "ok", http.StatusOK},
		{"unknown route", "GET", "/missing", "", "", http.StatusNotFound},
		{"unsupported method", "POST", "/health", "", "", http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			newHandler().ServeHTTP(response, httptest.NewRequest(tt.method, tt.path, nil))
			if response.Code != tt.status {
				t.Fatalf("status = %d, esperado %d", response.Code, tt.status)
			}
			if tt.key == "" {
				return
			}
			if got := response.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
				t.Fatalf("Content-Type inesperado: %q", got)
			}
			var body map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body[tt.key] != tt.value {
				t.Fatalf("resposta inesperada: %v", body)
			}
		})
	}
}
