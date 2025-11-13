package agent

import (
	"regexp"
	"strings"
)

// Redactor redacts PII from data structures.
type Redactor struct {
	blocklist []string
}

const redacted = "[REDACTED]"

var defaultBlocklist = []string{
	"password", "passwd", "pwd", "secret", "token",
	"api_key", "apikey", "credit_card", "card_number",
	"cvv", "ssn", "social_security", "pin", "passcode",
}

// NewRedactor creates a new redactor.
func NewRedactor(blocklist []string) *Redactor {
	if blocklist == nil {
		blocklist = defaultBlocklist
	}
	return &Redactor{blocklist: blocklist}
}

// Redact redacts PII from data.
func (r *Redactor) Redact(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		return r.redactMap(v)
	case []interface{}:
		return r.redactList(v)
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
			redacted[k] = redacted
		} else {
			redacted[k] = r.Redact(v)
		}
	}
	return redacted
}

func (r *Redactor) redactList(l []interface{}) []interface{} {
	redacted := make([]interface{}, len(l))
	for i, v := range l {
		redacted[i] = r.Redact(v)
	}
	return redacted
}

func (r *Redactor) redactString(s string) string {
	// Simple string redaction - could be enhanced with regex patterns
	return s
}

func (r *Redactor) shouldRedact(key string) bool {
	lowerKey := strings.ToLower(key)
	for _, pattern := range r.blocklist {
		if strings.Contains(pattern, "*") {
			regex := "^" + strings.ReplaceAll(regexp.QuoteMeta(pattern), "\\*", ".*") + "$"
			matched, _ := regexp.MatchString(regex, lowerKey)
			if matched {
				return true
			}
		} else if strings.Contains(lowerKey, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}
