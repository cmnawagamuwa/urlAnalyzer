package tests

import (
	"URLAnalyzer/checks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- DetectHTMLVersion ---

func TestDetectHTMLVersion(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "HTML5",
			html: "<!DOCTYPE html><html><body></body></html>",
			want: "HTML5",
		},
		{
			name: "HTML 4.01",
			html: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN">`,
			want: "HTML 4.01",
		},
		{
			name: "HTML 4.01 Strict",
			html: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Strict//EN">`,
			want: "HTML 4.01 Strict",
		},
		{
			name: "HTML 4.01 Transitional",
			html: `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">`,
			want: "HTML 4.01 Transitional",
		},
		{
			name: "XHTML 1.0 Strict",
			html: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">`,
			want: "XHTML 1.0 Strict",
		},
		{
			name: "XHTML 1.1",
			html: `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN" "http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd">`,
			want: "XHTML 1.1",
		},
		{
			name: "unknown DOCTYPE",
			html: "<!DOCTYPE something-custom><html></html>",
			want: "Unknown DOCTYPE",
		},
		{
			name: "no DOCTYPE",
			html: "<html><body><p>hello</p></body></html>",
			want: "No DOCTYPE found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checks.DetectHTMLVersion(tt.html)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// --- CountHeadings ---

func TestCountHeadings(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "no headings",
			html: "<html><body><p>hello</p></body></html>",
			want: "Headings: None found",
		},
		{
			name: "single H1",
			html: "<html><body><h1>Title</h1></body></html>",
			want: "Headings: H1=1, H2=0, H3=0, H4=0, H5=0, H6=0",
		},
		{
			name: "multiple mixed headings",
			html: "<h1>A</h1><h2>B</h2><h2>C</h2><h3>D</h3>",
			want: "Headings: H1=1, H2=2, H3=1, H4=0, H5=0, H6=0",
		},
		{
			name: "uppercase tags are counted",
			html: "<H1>Title</H1><H2>Sub</H2>",
			want: "Headings: H1=1, H2=1, H3=0, H4=0, H5=0, H6=0",
		},
		{
			name: "all heading levels",
			html: "<h1>1</h1><h2>2</h2><h3>3</h3><h4>4</h4><h5>5</h5><h6>6</h6>",
			want: "Headings: H1=1, H2=1, H3=1, H4=1, H5=1, H6=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checks.CountHeadings(tt.html)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// --- HasLoginForm ---

func TestHasLoginForm(t *testing.T) {
	tests := []struct {
		name string
		url  string
		html string
		want string
	}{
		{
			name: "login in URL",
			url:  "https://example.com/login",
			html: "",
			want: "Login form: Yes",
		},
		{
			name: "signin in URL",
			url:  "https://example.com/signin",
			html: "",
			want: "Login form: Yes",
		},
		{
			name: "password and form in HTML",
			url:  "https://example.com",
			html: `<form><input type="password" name="pass"><button>Submit</button></form>`,
			want: "Login form: Yes",
		},
		{
			name: "password and login keyword in HTML",
			url:  "https://example.com",
			html: `<div class="login"><input type="password"></div>`,
			want: "Login form: Yes",
		},
		{
			name: "no login form",
			url:  "https://example.com",
			html: "<html><body><p>Welcome to our site</p></body></html>",
			want: "Login form: No",
		},
		{
			name: "password alone is not enough",
			url:  "https://example.com",
			html: "<p>Use a strong password for security.</p>",
			want: "Login form: No",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checks.HasLoginForm(tt.url, tt.html)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// --- CountLinks ---

func TestCountLinks_InternalOnly(t *testing.T) {
	html := `<a href="/about">About</a><a href="/contact">Contact</a><a href="#">skip</a>`

	got, err := checks.CountLinks("https://example.com", html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(got, "Internal=<b>2</b>") {
		t.Errorf("expected 2 internal links, got: %q", got)
	}
	if !strings.Contains(got, "External=<b>0</b>") {
		t.Errorf("expected 0 external links, got: %q", got)
	}
}

// --- CheckReachable ---

func TestCheckReachable(t *testing.T) {
	t.Run("returns success for 200 response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		got, err := checks.CheckReachable(server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.Contains(got, "200") {
			t.Errorf("expected '200' in result, got %q", got)
		}
	})

	t.Run("returns error for non-200 response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		_, err := checks.CheckReachable(server.URL)
		if err == nil {
			t.Error("expected an error for non-200 response, got nil")
		}
	})

	t.Run("returns error for unreachable URL", func(t *testing.T) {
		_, err := checks.CheckReachable("http://localhost:0")
		if err == nil {
			t.Error("expected an error for unreachable URL, got nil")
		}
	})
}
