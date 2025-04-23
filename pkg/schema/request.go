package schema

// MCPRequest represents a basic JSON-RPC / MCP-compatible request
// used for LLM tool calls or context-based evaluations.
type MCPRequest struct {
	ID      string                 `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	Context map[string]interface{} `json:"context"`
}

// Validate checks if the request has required fields and formats
func Validate(req *MCPRequest) error {
	if req.ID == "" || req.Method == "" {
		return ErrInvalidRequest
	}
	return nil
}

// ErrInvalidRequest is returned when a request is malformed
var ErrInvalidRequest = &ValidationError{"missing required fields: id or method"}

// ValidationError represents a schema validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
