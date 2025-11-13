package agent

import (
	"regexp"
	"strings"
)

var defaultBlocklist = []string{
	"password", "passwd", "pwd", "secret", "token", "api_key", "apikey",
	"credit_card", "card_number", "cvv", "ssn", "social_security", "pin", "passcode",
}

// Redactor redacts PII from data structures.
type Redactor struct {
	patterns []*regexp.Regexp
}

// NewRedactor creates a new Redactor with the given rules.
func NewRedactor(rules []RedactionRule) *Redactor {
	patterns := make([]*regexp.Regexp, 0)

	// Add default patterns
	for _, field := range defaultBlocklist {
		pattern := regexp.MustCompile("(?i).*" + regexp.QuoteMeta(field) + ".*")
		patterns = append(patterns, pattern)
	}

	// Add custom rules
	for _, rule := range rules {
		if rule.Pattern != "" {
			pattern, err := regexp.Compile("(?i)" + rule.Pattern)
			if err == nil {
				patterns = append(patterns, pattern)
			}
		}
	}

	return &Redactor{patterns: patterns}
}

// Redact redacts PII from data.
func (r *Redactor) Redact(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	switch v := data.(type) {
	case map[string]interface{}:
		return r.redactMap(v)
	case []interface{}:
		return r.redactSlice(v)
	case string:
		return r.redactString(v)
	default:
		return data
	}
}

func (r *Redactor) redactMap(m map[string]interface{}) map[string]interface{} {
	redacted := make(map[string]interface{})
	for k, v := range m {
		if r.shouldRedact(k) {
			redacted[k] = "***REDACTED***"
		} else {
			redacted[k] = r.Redact(v)
		}
	}
	return redacted
}

func (r *Redactor) redactSlice(s []interface{}) []interface{} {
	redacted := make([]interface{}, len(s))
	for i, v := range s {
		redacted[i] = r.Redact(v)
	}
	return redacted
}

func (r *Redactor) redactString(s string) string {
	// Simple string redaction - could be enhanced
	return s
}

func (r *Redactor) shouldRedact(fieldName string) bool {
	for _, pattern := range r.patterns {
		if pattern.MatchString(fieldName) {
			return true
		}
	}
	return false
}
