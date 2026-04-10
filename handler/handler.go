package handler

import (
	"URLAnalyzer/analyzer"
	"URLAnalyzer/domain"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
)

// Renders the analysis result as a full HTML page.
// Automatically escapes all values, preventing XSS.
var resultTmpl = template.Must(template.New("result").Parse(`<!DOCTYPE html>
<html>
<head>
	<title>URL Analyzer - Result</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 800px; margin: 40px auto; padding: 0 20px; background-color: #f5f5f5; color: #333; }
		h1 { color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px; }
		a { color: #3498db; }
		.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-top: 24px; }
		.card { background-color: #fff; border: 1px solid #ddd; border-radius: 6px; padding: 16px; }
		.card.full { grid-column: 1 / -1; }
		.label { display: block; font-size: 12px; color: #888; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 6px; }
		.value { display: block; font-size: 15px; color: #2c3e50; word-break: break-all; }
		pre.links { margin: 0; font-size: 13px; line-height: 1.7; white-space: pre-wrap; color: #2c3e50; }
	</style>
</head>
<body>
	<h1>URL Analyzer</h1>
	<a href="/">Analyze another URL</a>
	<div class="grid">
		<div class="card full"><span class="label">URL</span><span class="value">{{.URL}}</span></div>
		<div class="card"><span class="label">Status</span><span class="value">{{.Reachable}}</span></div>
		<div class="card"><span class="label">HTML Version</span><span class="value">{{.HTMLVersion}}</span></div>
		<div class="card"><span class="label">Headings</span><span class="value">{{.Headings}}</span></div>
		<div class="card"><span class="label">Login Form</span><span class="value">{{.LoginForm}}</span></div>
		<div class="card full"><span class="label">Links</span><pre class="links">{{.Links}}</pre></div>
	</div>
</body>
</html>`))

// Serves the URL input form.
func ShowForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./webPages/index.html")
}

// Validates the submitted URL, runs the analysis, and renders the result page.
func HandleForm(w http.ResponseWriter, r *http.Request) {
	rawURL := r.FormValue("url")

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		http.Error(w, "Invalid URL. Please enter a valid http or https URL.", http.StatusBadRequest)
		return
	}

	slog.Info("Analyzing URL", "url", rawURL)

	result, err := analyzer.ProcessURL(rawURL)
	if err != nil {
		slog.Error("Analysis failed", "url", rawURL, "err", err)
		http.Error(w, "Could not analyze URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Embed the result but override Links as template.HTML so <b> tags render correctly.
	// All other fields are still accessible via the embedded struct.
	if err := resultTmpl.Execute(w, struct {
		domain.AnalysisResult
		Links template.HTML
	}{result, template.HTML(result.Links)}); err != nil {
		slog.Error("Failed to render result page", "err", err)
	}
}
