#!/bin/bash
set -e

# Use the internal service DNS and port 80 (Service port)
GATEWAY_URL="http://api-gateway.payments-system.svc.cluster.local:80"

echo "Checking health..."
curl -s "$GATEWAY_URL/health/readiness" || echo "Health check failed"

echo -e "\n\n1. Creating Merchant..."
# Create a merchant to get an API Key
# Now using /v1/merchants
MERCHANT_RESP=$(curl -s -X POST "$GATEWAY_URL/v1/merchants" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Merchant", "email": "test@merchant.com"}')

echo "Response: $MERCHANT_RESP"
API_KEY=$(echo $MERCHANT_RESP | grep -o '"api_key":"[^"'] *' | cut -d'"' -f4)

if [ -z "$API_KEY" ]; then
  echo "Failed to get API Key"
  exit 1
fi

echo "API Key: $API_KEY"

echo -e "\n2. Processing Payment..."
# Now using /v1/payments
PAYMENT_RESP=$(curl -s -X POST "$GATEWAY_URL/v1/payments" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d '{ "amount": 100.50, "currency": "USD", "card_number": "4242424242424242", "expiry_month": 12, "expiry_year": 2030, "cvv": "123" }')

echo "Response: $PAYMENT_RESP"

if [[ "$PAYMENT_RESP" == *"payment_id"* ]]; then
    echo -e "\n✅ SUCCESS: Payment Processed!"
else
    echo -e "\n❌ FAILURE: Payment Failed"
    exit 1
fi