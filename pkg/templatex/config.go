package templatex

import (
	"errors"
	"strings"
	"time"
)

// Config holds the configuration for a client.
type Config struct {
	Name    string
	Timeout time.Duration
	Secret  string
}

// SanitizedConfig is a copy of Config with sensitive fields redacted.
type SanitizedConfig struct {
	Name    string
	Timeout time.Duration
	Secret  string
}

// sensitiveFieldNames lists field name substrings that should be redacted.
var sensitiveFieldNames = []string{"secret", "token", "password", "key", "credential", "dsn", "url"}

// Validate checks required fields and rejects invalid values.
func (c Config) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return NewError(KindValidation, "Config.Validate", "name is required", errors.New("name is required"))
	}
	if c.Timeout < 0 {
		return NewError(KindValidation, "Config.Validate", "timeout must not be negative", errors.New("timeout must not be negative"))
	}
	return nil
}

// Sanitize returns a copy with sensitive fields redacted.
func (c Config) Sanitize() SanitizedConfig {
	return SanitizedConfig{
		Name:    c.Name,
		Timeout: c.Timeout,
		Secret:  redact("secret", c.Secret),
	}
}

// redact replaces the value with "***" if the field name matches a sensitive pattern.
// Empty values are preserved as-is.
func redact(fieldName, value string) string {
	if value == "" {
		return ""
	}
	lower := strings.ToLower(fieldName)
	for _, pattern := range sensitiveFieldNames {
		if strings.Contains(lower, pattern) {
			return "***"
		}
	}
	return value
}
