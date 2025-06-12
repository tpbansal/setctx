#!/bin/bash

# Create testdata directory if it doesn't exist
mkdir -p testdata

# Start Redis if not running
if ! redis-cli ping > /dev/null 2>&1; then
    echo "Starting Redis..."
    redis-server --daemonize yes
fi

# Build and run the example
echo "Building and running example..."
go run cmd/example/main.go 