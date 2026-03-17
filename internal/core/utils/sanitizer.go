package utils

import "github.com/microcosm-cc/bluemonday"

func SanitizeText(input string) string {
	// This removes ALL HTML tags, leaving only plain text
	// We'll only accept markdown with no html tags or custom styling or scripts
	// as not negate any possible XSS
	p := bluemonday.StrictPolicy()
	return p.Sanitize(input)
}
