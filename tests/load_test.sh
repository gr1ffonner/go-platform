#!/bin/bash

# Aggressive load test for dogs endpoint
# Tests maximum load on GET /api/v1/dogs/borzoi/image

set -e

echo "🔥 Starting aggressive load test with Hey"
echo "Target: Dogs API endpoint"
echo ""

# Check if hey is installed
if ! command -v hey &> /dev/null; then
    echo "❌ Hey not found. Installing..."
    go install github.com/rakyll/hey@latest
    echo "✅ Hey installed!"
fi

# Check if app is running
if ! curl -s http://localhost:8080/live > /dev/null; then
    echo "❌ App is not running on localhost:8080"
    echo "Please start your app first with: make run-pg"
    exit 1
fi

echo "📊 Running progressive load test..."
echo "   - Total Requests: 10000"
echo "   - Ramp up: 0 to 100 concurrent users over 30 seconds"
echo "   - Target: /api/v1/dogs/borzoi/image endpoint"
echo ""

# Run the progressive test (ramp up from 0 to 100 concurrent users over 30 seconds)
hey -n 10000 -c 100 -z 30s http://localhost:8080/api/v1/dogs/borzoi/image

echo ""
echo "🎯 Test completed!"
