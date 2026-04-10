package checks

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// CountLinks categorises all links as internal, external, or broken and returns a summary.
func CountLinks(pageURL string, html string) (string, error) {
	base, _ := url.Parse(pageURL)
	domain := base.Hostname()

	var internalLinks []string
	var externalHrefs []string

	for _, link := range extractLinks(html) {
		if strings.HasPrefix(link, "/") || strings.Contains(link, domain) {
			internalLinks = append(internalLinks, link)
		} else if strings.HasPrefix(link, "http") {
			externalHrefs = append(externalHrefs, link)
		}
	}

	// Check all external links concurrently instead of one by one
	type checkedLink struct {
		display string
		broken  bool
	}

	resultsCh := make(chan checkedLink, len(externalHrefs))

	var wg sync.WaitGroup
	for _, link := range externalHrefs {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			code := getLinkStatusCode(l)
			resultsCh <- checkedLink{
				display: fmt.Sprintf("%s [%d]", l, code),
				broken:  code == 0 || code >= 400,
			}
		}(link)
	}

	// Close the channel once all goroutines finish so we can range over it
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	var externalLinks []string
	var brokenLinks []string

	for r := range resultsCh {
		externalLinks = append(externalLinks, r.display)
		if r.broken {
			brokenLinks = append(brokenLinks, r.display)
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Links: Internal=<b>%d</b>, External=<b>%d</b>, Broken=<b>%d</b>",
		len(internalLinks), len(externalLinks), len(brokenLinks))

	if len(internalLinks) > 0 {
		sb.WriteString("\n  Internal Links:")
		for _, link := range internalLinks {
			sb.WriteString("\n    - " + link)
		}
	}

	if len(externalLinks) > 0 {
		sb.WriteString("\n  External Links:")
		for _, link := range externalLinks {
			sb.WriteString("\n    - " + link)
		}
	}

	return sb.String(), nil
}

// extractLinks pulls all href values from the HTML, preserving original URL casing.
func extractLinks(html string) []string {
	lower := strings.ToLower(html)
	var links []string
	start := 0

	for {
		pos := strings.Index(lower[start:], "href=\"")
		if pos == -1 {
			break
		}
		pos += start + 6
		end := strings.Index(lower[pos:], "\"")
		if end == -1 {
			break
		}

		link := html[pos : pos+end]
		start = pos + end

		if link == "" || link == "/" || link == "." || strings.HasPrefix(link, "#") {
			continue
		}

		links = append(links, link)
	}

	return links
}

// getLinkStatusCode fetches the link and returns the HTTP status code (0 if unreachable).
func getLinkStatusCode(link string) int {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(link)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
