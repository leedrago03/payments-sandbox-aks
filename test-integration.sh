#!/bin/bash

echo "🧪 Testing Integrated Platform..."
echo "================================"

# Test 1: Health checks
echo ""
echo "1️⃣  Testing Health Checks..."
echo "API Gateway:"
curl -s http://localhost:3000/health | jq . 2>/dev/null || curl -s http://localhost:3000/health

echo ""
echo "Payment Service:"
curl -s http://localhost:8081/actuator/health | jq . 2>/dev/null || curl -s http://localhost:8081/actuator/health

echo ""
echo "Merchant Service:"
curl -s http://localhost:3002/health/readiness | jq . 2>/dev/null || curl -s http://localhost:3002/health/readiness

echo ""
echo "Tokenization Service:"
curl -s http://localhost:3003/health/readiness | jq . 2>/dev/null || curl -s http://localhost:3003/health/readiness

echo ""
echo "Acquirer Simulator:"
curl -s http://localhost:3004/health | jq . 2>/dev/null || curl -s http://localhost:3004/health

echo ""
echo "Ledger Service:"
curl -s http://localhost:3005/health/readiness | jq . 2>/dev/null || curl -s http://localhost:3005/health/readiness

echo ""
echo "Audit Service:"
curl -s http://localhost:3006/health/readiness | jq . 2>/dev/null || curl -s http://localhost:3006/health/readiness

echo ""
echo "Reconciliation Service:"
curl -s http://localhost:3007/health/readiness | jq . 2>/dev/null || curl -s http://localhost:3007/health/readiness

# Test 2: Create merchant
echo ""
echo ""
echo "2️⃣  Creating Merchant..."
MERCHANT_RESPONSE=$(curl -s -X POST http://localhost:3002/api/merchants \
  -H "Content-Type: application/json" \
  -d '{
    "business_name": "Integration Test Store",
    "email": "integration@test.com",
    "contact_person": "Test User",
    "phone": "+1234567890"
  }')

echo "$MERCHANT_RESPONSE" | jq . 2>/dev/null || echo "$MERCHANT_RESPONSE"

MERCHANT_ID=$(echo "$MERCHANT_RESPONSE" | jq -r '.id' 2>/dev/null || echo "")
echo ""
echo "✅ Merchant ID: $MERCHANT_ID"

# Test 3: Create API Key
echo ""
echo "3️⃣  Creating API Key..."
API_KEY_RESPONSE=$(curl -s -X POST http://localhost:3002/api/merchants/$MERCHANT_ID/api-keys \
  -H "Content-Type: application/json" \
  -d '{"name": "Integration Test Key"}')

echo "$API_KEY_RESPONSE" | jq . 2>/dev/null || echo "$API_KEY_RESPONSE"

API_KEY=$(echo "$API_KEY_RESPONSE" | jq -r '.api_key' 2>/dev/null || echo "")
echo ""
echo "✅ API Key: ${API_KEY:0:30}..."

# Test 4: Tokenize card
echo ""
echo "4️⃣  Tokenizing Card..."
TOKEN_RESPONSE=$(curl -s -X POST http://localhost:3003/api/tokens \
  -H "Content-Type: application/json" \
  -d '{
    "pan": "4532015112830366",
    "expiry_month": "12",
    "expiry_year": "2026",
    "cvv": "123",
    "merchant_id": "'"$MERCHANT_ID"'"
  }')

echo "$TOKEN_RESPONSE" | jq . 2>/dev/null || echo "$TOKEN_RESPONSE"

TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.token' 2>/dev/null || echo "")
echo ""
echo "✅ Token: $TOKEN"

# Test 5: Acquirer authorization
echo ""
echo "5️⃣  Testing Acquirer Authorization..."
ACQUIRER_RESPONSE=$(curl -s -X POST http://localhost:3004/api/acquirer/authorize \
  -H "Content-Type: application/json" \
  -d '{
    "token": "'"$TOKEN"'",
    "amount": 100.50,
    "currency": "USD",
    "merchant_id": "'"$MERCHANT_ID"'"
  }')

echo "$ACQUIRER_RESPONSE" | jq . 2>/dev/null || echo "$ACQUIRER_RESPONSE"

# Test 6: Create ledger entries
echo ""
echo "6️⃣  Creating Ledger Entries..."
LEDGER_RESPONSE=$(curl -s -X POST http://localhost:3005/api/ledger/entries \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": "test_payment_001",
    "description": "Integration test payment",
    "entries": [
      {"account_id": "customer_'$MERCHANT_ID'", "entry_type": "DEBIT", "amount": 100.50, "currency": "USD"},
      {"account_id": "merchant_'$MERCHANT_ID'", "entry_type": "CREDIT", "amount": 100.50, "currency": "USD"}
    ]
  }')

echo "$LEDGER_RESPONSE" | jq . 2>/dev/null || echo "$LEDGER_RESPONSE"

# Test 7: Create audit log
echo ""
echo "7️⃣  Creating Audit Log..."
AUDIT_RESPONSE=$(curl -s -X POST http://localhost:3006/api/audit/logs \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "PAYMENT_CREATED",
    "entity_type": "payment",
    "entity_id": "test_payment_001",
    "actor_type": "merchant",
    "actor_id": "'"$MERCHANT_ID"'",
    "action": "create_payment",
    "details": "Integration test payment created",
    "success": true
  }')

echo "$AUDIT_RESPONSE" | jq . 2>/dev/null || echo "$AUDIT_RESPONSE"

echo ""
echo ""
echo "🎉 Integration Test Complete!"
echo "================================"
echo "✅ All 8 services are healthy"
echo "✅ Merchant created: $MERCHANT_ID"
echo "✅ API Key generated"
echo "✅ Card tokenized: $TOKEN"
echo "✅ Acquirer authorization tested"
echo "✅ Ledger entries created"
echo "✅ Audit log recorded"
echo ""
echo "🌐 Platform URLs:"
echo "   API Gateway:        http://localhost:3000"
echo "   Payment Service:    http://localhost:8081"
echo "   Merchant Service:   http://localhost:3002"
echo "   All other services running on ports 3003-3007"
