package policy

import "safectx/pkg/schema"

// Evaluate a request against a custom policy (dummy logic for now)
func Evaluate(req *schema.MCPRequest) (bool, error) {
	if req.Method == "sensitive_method" {
		return false, nil // Deny request based on method
	}
	return true, nil
}
