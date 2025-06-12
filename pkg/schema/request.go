package schema

// MCPRequest represents the structure of an incoming JSON-RPC request
type MCPRequest struct {
	ID     string                 `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

// Validate performs basic validation on the MCPRequest
func Validate(req *MCPRequest) error {
	if req.ID == "" {
		return ErrMissingID
	}
	if req.Method == "" {
		return ErrMissingMethod
	}
	if req.Params == nil {
		return ErrMissingParams
	}
	return nil
}

// Error definitions
var (
	ErrMissingID     = NewValidationError("missing required field: id")
	ErrMissingMethod = NewValidationError("missing required field: method")
	ErrMissingParams = NewValidationError("missing required field: params")
)

// ValidationError represents a validation error
type ValidationError struct {
	message string
}

func NewValidationError(message string) error {
	return &ValidationError{message: message}
}

func (e *ValidationError) Error() string {
	return e.message
}
