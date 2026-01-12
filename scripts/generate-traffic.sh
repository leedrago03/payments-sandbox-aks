#!/bin/sh
# Script with random email to ensure fresh API keys
RAND=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
EMAIL="test-$RAND@enterprise.com"
M_URL="http://merchant-service.payments-system.svc.cluster.local:3002/api/merchants"
G_URL="http://api-gateway.payments-system.svc.cluster.local:80/v1/payments"

echo "🚀 Starting Fresh Stress Test with $EMAIL..."

# 1. Create Merchant
RESP=$(curl -s -X POST $M_URL -H "Content-Type: application/json" -d "{\"business_name\":\"Test-$RAND\",\"email\":\"$EMAIL\"}")
MID=$(echo $RESP | sed 's/.*"id":"\([^"]*\)".*/\1/')
echo "Merchant: $MID"

# 2. Create Key
KRESP=$(curl -s -X POST $M_URL/$MID/api-keys -H "Content-Type: application/json" -d '{"name":"Key"}')
AKEY=$(echo $KRESP | sed 's/.*"api_key":"\([^"]*\)".*/\1/')
echo "Key: OK"

# 3. Loop Payments
echo "💸 Sending 100 Payments..."
for i in $(seq 1 100); do
  CODE=$(curl -s -o /dev/null -w "% {http_code}" -X POST $G_URL \
    -H "X-API-Key: $AKEY" \
    -H "Content-Type: application/json" \
    -d "{\"amount\":10.$i,\"currency\":\"USD\",\"card_number\":\"4532015112830366\"}")
  
  if [ "$CODE" = "200" ] || [ "$CODE" = "201" ]; then
    echo -n "."
  else
    echo -n "X($CODE)"
  fi
  sleep 0.1
done
echo ""
echo "🎉 Done!"
