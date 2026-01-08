#!/bin/bash

echo "🚀 Starting Payments Platform..."
echo "================================"

# Build and start all services
docker-compose up --build -d

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 30

echo ""
echo "✅ Platform Status:"
echo "================================"
docker-compose ps

echo ""
echo "🌐 Service URLs:"
echo "================================"
echo "API Gateway:              http://localhost:3000"
echo "Payment Service:          http://localhost:8081"
echo "Merchant Service:         http://localhost:3002"
echo "Tokenization Service:     http://localhost:3003"
echo "Acquirer Simulator:       http://localhost:3004"
echo "Ledger Service:           http://localhost:3005"
echo "Audit Service:            http://localhost:3006"
echo "Reconciliation Service:   http://localhost:3007"

echo ""
echo "📊 View logs:"
echo "================================"
echo "docker-compose logs -f [service-name]"
echo "docker-compose logs -f payment-service"

echo ""
echo "🛑 Stop platform:"
echo "================================"
echo "docker-compose down"
