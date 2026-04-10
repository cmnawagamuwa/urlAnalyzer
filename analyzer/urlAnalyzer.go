package analyzer

import (
	"URLAnalyzer/checks"
	"URLAnalyzer/domain"
	"URLAnalyzer/util"
	"log/slog"
)

// ProcessURL runs all analysis checks on the given URL and returns the result.
func ProcessURL(rawURL string) (domain.AnalysisResult, error) {
	result := domain.AnalysisResult{URL: rawURL}

	reachable, err := checks.CheckReachable(rawURL)
	if err != nil {
		return result, err
	}
	result.Reachable = reachable

	html, err := util.FetchHTML(rawURL)
	if err != nil {
		return result, err
	}

	// Each check runs in its own goroutine and sends its result back via a channel.
	// This way all four checks run at the same time instead of one after another.
	htmlVersionCh := make(chan string, 1)
	headingsCh    := make(chan string, 1)
	linksCh       := make(chan string, 1)
	loginFormCh   := make(chan string, 1)

	go func() {
		v, err := checks.DetectHTMLVersion(html)
		if err != nil {
			slog.Warn("DetectHTMLVersion failed", "err", err)
		}
		htmlVersionCh <- v
	}()

	go func() {
		v, err := checks.CountHeadings(html)
		if err != nil {
			slog.Warn("CountHeadings failed", "err", err)
		}
		headingsCh <- v
	}()

	go func() {
		v, err := checks.CountLinks(rawURL, html)
		if err != nil {
			slog.Warn("CountLinks failed", "err", err)
		}
		linksCh <- v
	}()

	go func() {
		v, err := checks.HasLoginForm(rawURL, html)
		if err != nil {
			slog.Warn("HasLoginForm failed", "err", err)
		}
		loginFormCh <- v
	}()

	// Wait for all four results
	result.HTMLVersion = <-htmlVersionCh
	result.Headings    = <-headingsCh
	result.Links       = <-linksCh
	result.LoginForm   = <-loginFormCh

	return result, nil
}
