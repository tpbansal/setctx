package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"safectx/internal/middleware"
	"safectx/internal/rpc"
)

func main() {
	// Initialize the HTTP request multiplexer
	mux := http.NewServeMux()

	// Create the handler wrapped with middleware chain
	handler := middleware.WithRecovery(
		middleware.WithLogging(
			middleware.WithValidation(
				middleware.WithPolicy(
					middleware.WithRedaction(http.HandlerFunc(rpc.HandleRequest)),
				),
			),
		),
	)

	// Register the route with the handler
	mux.Handle("/rpc", handler)

	// Initialize the HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a goroutine so it doesn't block
	go func() {
		log.Println("SafeCtx MCP gateway listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down the server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
