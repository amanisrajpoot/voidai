package agent

import (
	"regexp"
	"strings"
)

var defaultPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[:=]\s*([^\s,}]+)`),
	regexp.MustCompile(`(?i)(api[_-]?key|apikey)\s*[:=]\s*([^\s,}]+)`),
	regexp.MustCompile(`(?i)(token|bearer)\s*[:=]\s*([^\s,}]+)`),
	regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`), // Credit card
	regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`),                      // SSN
}

// Redaction handles PII redaction
type Redaction struct {
	patterns []*regexp.Regexp
}

// NewRedaction creates a new redaction instance
func NewRedaction(customRules []string) *Redaction {
	patterns := make([]*regexp.Regexp, len(defaultPatterns))
	copy(patterns, defaultPatterns)

	if customRules != nil {
		for _, rule := range customRules {
			if re, err := regexp.Compile(rule); err == nil {
				patterns = append(patterns, re)
			}
		}
	}

	return &Redaction{patterns: patterns}
}

// Redact redacts sensitive information from a string
func (r *Redaction) Redact(input string) string {
	if input == "" {
		return input
	}

	result := input
	for _, pattern := range r.patterns {
		result = pattern.ReplaceAllString(result, "[REDACTED]")
	}
	return result
}

// RedactMap redacts sensitive information from a map
func (r *Redaction) RedactMap(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return make(map[string]interface{})
	}

	redacted := make(map[string]interface{})
	for key, value := range data {
		if shouldRedactKey(key) {
			redacted[key] = "[REDACTED]"
		} else if str, ok := value.(string); ok {
			redacted[key] = r.Redact(str)
		} else if m, ok := value.(map[string]interface{}); ok {
			redacted[key] = r.RedactMap(m)
		} else {
			redacted[key] = value
		}
	}
	return redacted
}

func shouldRedactKey(key string) bool {
	lowerKey := strings.ToLower(key)
	return strings.Contains(lowerKey, "password") ||
		strings.Contains(lowerKey, "passwd") ||
		strings.Contains(lowerKey, "pwd") ||
		strings.Contains(lowerKey, "api_key") ||
		strings.Contains(lowerKey, "apikey") ||
		strings.Contains(lowerKey, "token") ||
		strings.Contains(lowerKey, "secret")
}
