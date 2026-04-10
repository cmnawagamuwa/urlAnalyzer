package tests

import (
	"URLAnalyzer/handler"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestHandleForm_InvalidURL checks that bad URLs are rejected before any HTTP call is made.
func TestHandleForm_InvalidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"empty string", ""},
		{"no scheme", "example.com"},
		{"ftp scheme is not allowed", "ftp://example.com"},
		{"plain text", "not-a-url"},
		{"missing host", "http://"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{"url": {tt.url}}
			req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			rec := httptest.NewRecorder()
			handler.HandleForm(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("got status %d, want %d (400 Bad Request)", rec.Code, http.StatusBadRequest)
			}
		})
	}
}
