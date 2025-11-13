package agent

import (
	"reflect"
	"regexp"
	"strings"
)

const defaultReplacement = "***REDACTED***"

// Redactor redacts PII from data structures.
type Redactor struct {
	rules []RedactionRule
}

// NewRedactor creates a new redactor with the given rules.
func NewRedactor(rules []RedactionRule) *Redactor {
	return &Redactor{rules: rules}
}

// Redact redacts PII from the given data.
func (r *Redactor) Redact(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Map:
		return r.redactMap(data)
	case reflect.Slice, reflect.Array:
		return r.redactSlice(data)
	case reflect.String:
		return r.redactString(data.(string))
	default:
		return data
	}
}

func (r *Redactor) redactMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(data)

	for _, key := range val.MapKeys() {
		keyStr := key.String()
		value := val.MapIndex(key).Interface()

		shouldRedact := false
		for _, rule := range r.rules {
			if rule.Field != "" && strings.EqualFold(rule.Field, keyStr) {
				shouldRedact = true
				break
			}
			if rule.Pattern != "" {
				matched, err := regexp.MatchString("(?i)"+rule.Pattern, keyStr)
				if err == nil && matched {
					shouldRedact = true
					break
				}
			}
		}

		if shouldRedact {
			replacement := defaultReplacement
			for _, rule := range r.rules {
				if rule.Field == keyStr || (rule.Pattern != "" && r.matchesPattern(rule.Pattern, keyStr)) {
					if rule.Replacement != "" {
						replacement = rule.Replacement
					}
					break
				}
			}
			result[keyStr] = replacement
		} else {
			result[keyStr] = r.Redact(value)
		}
	}

	return result
}

func (r *Redactor) redactSlice(data interface{}) []interface{} {
	val := reflect.ValueOf(data)
	result := make([]interface{}, val.Len())

	for i := 0; i < val.Len(); i++ {
		result[i] = r.Redact(val.Index(i).Interface())
	}

	return result
}

func (r *Redactor) redactString(str string) string {
	result := str
	for _, rule := range r.rules {
		if rule.Pattern != "" {
			replacement := rule.Replacement
			if replacement == "" {
				replacement = defaultReplacement
			}
			re := regexp.MustCompile("(?i)" + rule.Pattern)
			result = re.ReplaceAllString(result, replacement)
		}
	}
	return result
}

func (r *Redactor) matchesPattern(pattern, str string) bool {
	matched, err := regexp.MatchString("(?i)"+pattern, str)
	return err == nil && matched
}
