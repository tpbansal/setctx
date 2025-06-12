package policy

import "safectx/pkg/schema"

// Engine defines the interface for policy evaluation
type Engine interface {
	// Evaluate checks if a request is allowed based on the policy rules
	Evaluate(req *schema.MCPRequest) (bool, error)
}

// DefaultEngine implements the Engine interface
type DefaultEngine struct {
	// Add any necessary fields here
}

// NewDefaultEngine creates a new instance of DefaultEngine
func NewDefaultEngine() *DefaultEngine {
	return &DefaultEngine{}
}

// Evaluate implements the Engine interface
func (e *DefaultEngine) Evaluate(req *schema.MCPRequest) (bool, error) {
	// TODO: Implement policy evaluation logic
	return true, nil
}
