#!/bin/bash

# Quick load test using Hey
# Faster alternative to Vegeta for simple testing

set -e

echo "🚀 Starting quick load test with Hey"
echo "Target: Go Platform API"
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

echo "📊 Running quick test..."
echo "   - Requests: 1000"
echo "   - Concurrency: 10"
echo "   - Target: /live endpoint"
echo ""

# Run the test
hey -n 1000 -c 10 http://localhost:8080/live

