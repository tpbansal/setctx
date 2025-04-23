package rpc

import (
	"encoding/json"
	"net/http"
	"safectx/internal/contextfilter"
	"safectx/internal/policy"

	// "safectx/internal/rpc/schema"
	"safectx/pkg/schema"
)

// NewGatewayHandler returns the main SafeCtx HTTP handler
func NewGatewayHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req schema.MCPRequest

		// Decode request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Validate schema
		if err := schema.Validate(&req); err != nil {
			http.Error(w, "Schema validation failed", http.StatusBadRequest)
			return
		}

		// Policy check
		allowed, err := policy.Evaluate(&req)
		if err != nil || !allowed {
			http.Error(w, "Policy denied request", http.StatusForbidden)
			return
		}

		// Redact sensitive content
		contextfilter.Redact(&req)

		// Forward (mocked for now)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Request processed by SafeCtx",
		})
	})
}
