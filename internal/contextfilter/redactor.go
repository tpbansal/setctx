package contextfilter

import (
	"safectx/pkg/schema"
)

// Redact removes sensitive information from the request
func Redact(req *schema.MCPRequest) {
	if req.Params == nil {
		return
	}

	// List of sensitive keys to redact
	sensitiveKeys := []string{
		"password",
		"api_key",
		"secret",
		"token",
		"credentials",
	}

	// Redact sensitive values
	for _, key := range sensitiveKeys {
		if _, exists := req.Params[key]; exists {
			req.Params[key] = "[REDACTED]"
		}
	}
}
