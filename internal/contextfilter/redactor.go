package contextfilter

import (
	"fmt"
	"safectx/pkg/schema"
)

// Redact will scan the provided MCPRequest and sanitize sensitive fields.
func Redact(req *schema.MCPRequest) {
	// List of sensitive fields we want to redact
	sensitiveFields := []string{"password", "api_key", "ssn", "credit_card", "secret_key"}

	// Iterate over all params and redact sensitive ones
	for _, field := range sensitiveFields {
		if value, exists := req.Params[field]; exists {
			req.Params[field] = fmt.Sprintf("[REDACTED: %s]", value)
		}
	}

	// Optionally redact context or other parts of the request
	if req.Context != nil {
		// Example: Redact any sensitive keys in the context
		if _, exists := req.Context["user_password"]; exists {
			req.Context["user_password"] = "[REDACTED]"
		}
	}
}
