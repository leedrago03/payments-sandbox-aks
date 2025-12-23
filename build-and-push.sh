#!/bin/bash
set -e

# Configuration
ACR_NAME="paymentssandboxacr" # Replace with your actual ACR name if known, or pass as arg
TAG="latest"

if [ -z "$1" ]; then
  echo "Usage: ./build-and-push.sh <ACR_NAME> [TAG]"
  echo "Defaulting to ACR: $ACR_NAME"
else
  ACR_NAME=$1
fi

if [ -n "$2" ]; then
  TAG=$2
fi

echo "Logging into ACR $ACR_NAME..."
az acr login --name $ACR_NAME

echo "Building and Pushing Services..."

services=(
  "payment-service"
  "tokenization-service"
  "audit-service"
  "ledger-service"
  "merchant-service"
  "api-gateway"
  "acquirer-simulator"
  "reconciliation-service"
)

for service in "${services[@]}"; do
  echo "------------------------------------------------"
  echo "Processing $service..."
  
  # Build using the root context logic we established
  docker build -f services/$service/Dockerfile -t $ACR_NAME.azurecr.io/payments-sandbox/$service:$TAG .
  
  echo "Pushing $service..."
  docker push $ACR_NAME.azurecr.io/payments-sandbox/$service:$TAG
done

echo "------------------------------------------------"
echo "All services built and pushed successfully!"
