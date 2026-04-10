package checks

import "strings"

// HasLoginForm checks whether the page contains a login form.
func HasLoginForm(url string, html string) (string, error) {
	if strings.Contains(url, "login") || strings.Contains(url, "signin") {
		return "Login form: Yes", nil
	}

	lower := strings.ToLower(html)
	if strings.Contains(lower, "password") &&
		(strings.Contains(lower, "login") ||
			strings.Contains(lower, "signin") ||
			strings.Contains(lower, "form")) {
		return "Login form: Yes", nil
	}

	return "Login form: No", nil
}
