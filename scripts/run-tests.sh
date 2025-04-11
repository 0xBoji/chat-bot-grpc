#!/bin/bash

# Script to run tests for the project

set -e

# Print header
echo "Running tests for gRPC Messenger Core"
echo "===================================="

# Run tests for internal packages
echo "Running tests for internal packages..."
go test -v ./internal/...

# Run tests for db packages
echo "Running tests for db packages..."
go test -v ./db/...

# Print summary
echo "===================================="
echo "All tests completed successfully!"
