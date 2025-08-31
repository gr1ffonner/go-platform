#!/bin/bash

# Simple script to generate 100 health check requests
# This will help populate your metrics dashboard

echo "🚀 Starting 100 health check requests to http://localhost:8080/live"
echo "This will generate HTTP metrics for your Grafana dashboard"
echo ""

# Counter for tracking progress
count=0
total=100

# Loop to make 100 requests
for i in {1..100}; do
    # Make the request (silent output)
    curl -s http://localhost:8080/live > /dev/null
    
    # Increment counter
    ((count++))
    
    # Show progress every 10 requests
    if [ $((count % 10)) -eq 0 ]; then
        echo "✅ Completed $count/$total requests"
    fi
    
    # Small delay to spread requests over time
    sleep 0.1
done

echo ""
echo "🎉 All 100 requests completed!"
echo ""
echo "📊 Check your metrics:"
echo "   - Metrics endpoint: http://localhost:9090"
echo "   - Prometheus: http://localhost:9091"
echo "   - Grafana dashboard: http://localhost:3000"
echo ""
echo "💡 Try running this script multiple times to see more data in your dashboard!"
