#!/bin/bash
set -e

GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Configuring Ingress for Observability Dashboards...${NC}"

# 1. Get Ingress IP
echo "Waiting for Istio Ingress External IP..."
INGRESS_IP=""
while [ -z "$INGRESS_IP" ]; do
  az aks command invoke --resource-group rg-payments-aks --name payments-aks-dev --command "kubectl get svc istio-ingress -n istio-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}'" -o json > ingress_out.json
  
  # Extract logs using python
  INGRESS_IP=$(python3 -c "import json; print(json.load(open('ingress_out.json'))['logs'].strip())")
  
  if [ -z "$INGRESS_IP" ] || [ "$INGRESS_IP" == "None" ]; then
    echo "Waiting for IP... (Retrying in 10s)"
    INGRESS_IP=""
    sleep 10
  fi
done

rm ingress_out.json
echo -e "${GREEN}Found Ingress IP: $INGRESS_IP${NC}"

# 2. Update Manifest
MANIFEST_PATH="k8s-manifests/infrastructure/istio/dashboard-gateway.yaml"

# Replace placeholder or existing IP pattern
# We use a temp file to avoid sed issues
sed "s/INGRESS_IP_PLACEHOLDER/$INGRESS_IP/g" $MANIFEST_PATH > $MANIFEST_PATH.tmp && mv $MANIFEST_PATH.tmp $MANIFEST_PATH
# Also handle update scenario
sed -E "s/[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+\.nip\.io/$INGRESS_IP.nip.io/g" $MANIFEST_PATH > $MANIFEST_PATH.tmp && mv $MANIFEST_PATH.tmp $MANIFEST_PATH

echo "Updated $MANIFEST_PATH with $INGRESS_IP"

# 3. Commit and Push
echo -e "${GREEN}Pushing configuration to Git...${NC}"
git add $MANIFEST_PATH
git commit -m "config(ingress): update dashboard gateway with IP $INGRESS_IP" || echo "No changes to commit."
git push origin main

echo -e "${GREEN}Ingress Configuration Complete!${NC}"
echo "Dashboards will be available at:"
echo "- http://kiali.$INGRESS_IP.nip.io"
echo "- http://grafana.$INGRESS_IP.nip.io"
echo "- http://jaeger.$INGRESS_IP.nip.io"
echo "- http://argocd.$INGRESS_IP.nip.io"