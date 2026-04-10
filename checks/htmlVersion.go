package checks

import "strings"

// DetectHTMLVersion inspects the DOCTYPE declaration and returns the HTML version.
func DetectHTMLVersion(html string) (string, error) {
	doctype := strings.ToUpper(html[:min(512, len(html))])

	switch {
	case strings.HasPrefix(doctype, "<!DOCTYPE HTML>"):
		return "HTML5", nil
	case strings.Contains(doctype, "XHTML 1.1"):
		return "XHTML 1.1", nil
	case strings.Contains(doctype, "XHTML 1.0 STRICT"):
		return "XHTML 1.0 Strict", nil
	case strings.Contains(doctype, "XHTML 1.0 TRANSITIONAL"):
		return "XHTML 1.0 Transitional", nil
	case strings.Contains(doctype, "XHTML 1.0 FRAMESET"):
		return "XHTML 1.0 Frameset", nil
	case strings.Contains(doctype, "HTML 4.01") && strings.Contains(doctype, "STRICT"):
		return "HTML 4.01 Strict", nil
	case strings.Contains(doctype, "HTML 4.01") && strings.Contains(doctype, "TRANSITIONAL"):
		return "HTML 4.01 Transitional", nil
	case strings.Contains(doctype, "HTML 4.01") && strings.Contains(doctype, "FRAMESET"):
		return "HTML 4.01 Frameset", nil
	case strings.Contains(doctype, "HTML 4.01"):
		return "HTML 4.01", nil
	case strings.Contains(doctype, "HTML 3.2"):
		return "HTML 3.2", nil
	case strings.Contains(doctype, "HTML 2.0"):
		return "HTML 2.0", nil
	case strings.HasPrefix(doctype, "<!DOCTYPE"):
		return "Unknown DOCTYPE", nil
	default:
		return "No DOCTYPE found", nil
	}
}
