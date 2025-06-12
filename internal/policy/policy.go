package policy

import (
	"safectx/pkg/schema"
)

// Evaluate checks if the request is allowed based on the policy rules
func Evaluate(req *schema.MCPRequest) (bool, error) {
	// TODO: Implement actual policy evaluation logic
	// For now, we'll allow all requests
	return true, nil
}
