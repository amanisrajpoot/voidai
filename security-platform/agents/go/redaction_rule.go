package agent

// RedactionRule represents a redaction rule configuration.
type RedactionRule struct {
	Pattern     string `yaml:"pattern" json:"pattern"`
	Replacement string `yaml:"replacement" json:"replacement"`
	Field       string `yaml:"field" json:"field"`
}
