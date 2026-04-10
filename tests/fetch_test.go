package tests

import (
	"URLAnalyzer/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestFetchHTML uses a local test server to avoid real HTTP requests.
func TestFetchHTML(t *testing.T) {
	t.Run("returns page content on success", func(t *testing.T) {
		body := "<html><body>Hello</body></html>"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body))
		}))
		defer server.Close()

		got, err := util.FetchHTML(server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != body {
			t.Errorf("got %q, want %q", got, body)
		}
	})

	t.Run("returns error for unreachable URL", func(t *testing.T) {
		_, err := util.FetchHTML("http://localhost:0")
		if err == nil {
			t.Error("expected an error for unreachable URL, got nil")
		}
	})

	t.Run("preserves original casing of content", func(t *testing.T) {
		body := "<H1>TITLE</H1>"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(body))
		}))
		defer server.Close()

		got, err := util.FetchHTML(server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != body {
			t.Errorf("content casing changed: got %q, want %q", got, body)
		}
	})
}
