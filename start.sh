#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

mkdir -p data

echo "Starting Location Tracking Shortlink Server..."
echo "Admin panel: http://localhost:8080/admin"
echo "API endpoint: http://localhost:8080/api/shortlinks/create"

./track
