#!/bin/bash

# Deployment script for all services

set -e

# Build all services
./scripts/build.sh

# Create deployment directory
mkdir -p deploy

# Copy binaries
cp bin/auth-service deploy/
cp bin/chat-service deploy/
cp bin/room-service deploy/
cp bin/gateway deploy/

# Copy configuration
cp -r config deploy/

echo "Deployment package created in deploy/ directory"
