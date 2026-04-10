package checks

import (
	"fmt"
	"strings"
)

// CountHeadings counts H1–H6 tags in the HTML and returns a summary.
func CountHeadings(html string) (string, error) {
	lowerCasedHtml := strings.ToLower(html)
	h1 := strings.Count(lowerCasedHtml, "<h1")
	h2 := strings.Count(lowerCasedHtml, "<h2")
	h3 := strings.Count(lowerCasedHtml, "<h3")
	h4 := strings.Count(lowerCasedHtml, "<h4")
	h5 := strings.Count(lowerCasedHtml, "<h5")
	h6 := strings.Count(lowerCasedHtml, "<h6")

	if h1+h2+h3+h4+h5+h6 == 0 {
		return "Headings: None found", nil
	}

	return fmt.Sprintf("Headings: H1=%d, H2=%d, H3=%d, H4=%d, H5=%d, H6=%d",
		h1, h2, h3, h4, h5, h6), nil
}
