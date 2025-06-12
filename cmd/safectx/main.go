package main

import (
	"log"
	"net/http"
	"safectx/internal/middleware"
	"safectx/internal/rpc"
	"time"
)

func main() {
	// Create OIDC authenticator
	oidcAuth, err := middleware.NewOIDCAuthenticator(
		"https://your-oidc-provider",
		"client-id",
		"client-secret",
		"http://localhost:8080/callback",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create middleware chain
	chain := middleware.Chain(
		middleware.LoggingMiddleware(),
		middleware.RateLimitMiddleware(middleware.NewRateLimiter(10.0, 10.0, time.Second)),
		middleware.AuthMiddleware(oidcAuth),
	)
	// Create the gateway handler
	handler := rpc.NewGatewayHandler()

	// Start the HTTP server
	log.Println("Starting SafeCtx server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
