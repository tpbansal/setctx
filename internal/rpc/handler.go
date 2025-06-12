package rpc

import (
	"encoding/json"
	"log"
	"net/http"
	"safectx/internal/contextfilter"
	"safectx/internal/detection"
	"safectx/internal/policy"
	"safectx/pkg/schema"
)

// NewGatewayHandler returns the main SafeCtx HTTP handler
func NewGatewayHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req schema.MCPRequest

		// Decode the incoming request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			log.Printf("Error decoding request: %v", err)
			return
		}

		// Validate schema
		if err := schema.Validate(&req); err != nil {
			http.Error(w, "Schema validation failed: "+err.Error(), http.StatusBadRequest)
			log.Printf("Schema validation failed: %v", err)
			return
		}

		// Check for prompt injection
		if detection.CheckForInjection(&req) {
			http.Error(w, "Potential prompt injection detected", http.StatusForbidden)
			log.Printf("Prompt injection detected in request: %+v", req)
			return
		}

		// Evaluate policy
		policyEngine := policy.NewDefaultEngine()
		allowed, err := policyEngine.Evaluate(&req)
		if err != nil || !allowed {
			http.Error(w, "Policy denied request", http.StatusForbidden)
			log.Printf("Policy denied request: %v", err)
			return
		}

		// Redact sensitive content
		contextfilter.Redact(&req)

		// Log successful request processing
		log.Printf("Successfully processed request: %+v", req)

		// Respond with a success message
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status":  "ok",
			"message": "Request processed by SafeCtx",
			"data":    req.Params,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding response: %v", err)
		}
	})
}
