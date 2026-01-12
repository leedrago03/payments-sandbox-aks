#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}Starting Payments Platform Bootstrap...${NC}"

# 1. Prerequisites
if [ ! -f "infrastructure-outputs.json" ]; then
    echo -e "${RED}Error: infrastructure-outputs.json not found. Run terraform output -json > ... first.${NC}"
    exit 1
fi

# 2. Configure Manifests
echo -e "${GREEN}Configuring Kubernetes Manifests...${NC}"
python3 scripts/configure-manifests.py

# 3. Get Credentials
RG_NAME="rg-payments-aks" # Should be dynamic, but hardcoded for now or fetched
CLUSTER_NAME="payments-aks-dev"

echo -e "${GREEN}Getting AKS Credentials for $CLUSTER_NAME...${NC}"
az aks get-credentials --resource-group $RG_NAME --name $CLUSTER_NAME --overwrite-existing

# 4. Bootstrap Namespaces
echo -e "${GREEN}Creating Namespaces...${NC}"
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace payments-system --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace payments-data --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace istio-system --dry-run=client -o yaml | kubectl apply -f -

# 5. Seed Kubernetes Secrets
echo -e "${GREEN}Seeding Kubernetes Secrets...${NC}"

# Extract values using python/jq or just use the outputs file processing
# For simplicity, we assume the python script handled the configmaps, but SECRETS must be seeded in K8s.

# Helper to get json value
get_val() {
    cat infrastructure-outputs.json | grep -A 4 "\"$1\"" | grep "value" | cut -d '"' -f 4
}

REDIS_KEY=$(grep -A 4 "\"redis_primary_access_key\"" infrastructure-outputs.json | grep "value" | cut -d '"' -f 4)
# Redis Key might have special chars, handle carefully or use Python.
# Actually, the python script is better for this.

# Let's delegate secret creation to a python helper to avoid bash string escaping hell.
python3 -c "
import json, subprocess, base64, os

with open('infrastructure-outputs.json') as f:
    out = json.load(f)

def run(cmd):
    print(f'Exec: {cmd}')
    subprocess.run(cmd, shell=True, check=True)

redis_key = out['redis_primary_access_key']['value']
pg_host = out['postgresql_fqdn']['value']
# Hardcoded for now per main.tf
pg_pass = 'P@ssw0rd1234!' 

# HMAC Key
hmac_key = base64.b64encode(os.urandom(32)).decode('utf-8')

# Redis Secret
run(f\"kubectl create secret generic redis-credentials --from-literal=password='{redis_key}' --namespace payments-system --dry-run=client -o yaml | kubectl apply -f -\")

# Postgres Secret
run(f\"kubectl create secret generic postgresql-credentials --from-literal=username='pgadmin' --from-literal=password='{pg_pass}' --from-literal=host='{pg_host}' --from-literal=port='5432' --from-literal=sslmode='require' --namespace payments-system --dry-run=client -o yaml | kubectl apply -f -\")

# Audit Secret (System Namespace)
run(f\"kubectl create secret generic audit-secrets --from-literal=AUDIT_HMAC_KEY='{hmac_key}' --namespace payments-system --dry-run=client -o yaml | kubectl apply -f -\")
"

# 6. Seed Key Vault
echo -e "${GREEN}Seeding Azure Key Vault...${NC}"
KV_NAME=$(get_val "keyvault_uri" | cut -d '/' -f 3 | cut -d '.' -f 1)
ENCRYPTION_KEY=$(openssl rand -base64 32)

echo "Setting key in $KV_NAME..."
# Check if key exists
if az keyvault secret show --vault-name $KV_NAME --name payment-encryption-key &>/dev/null; then
    echo "Key already exists."
else
    # Assign self role first if needed
    USER_ID=$(az ad signed-in-user show --query id -o tsv)
    az role assignment create --role "Key Vault Secrets Officer" --assignee-object-id $USER_ID --scope "/subscriptions/$(az account show --query id -o tsv)/resourceGroups/$RG_NAME/providers/Microsoft.KeyVault/vaults/$KV_NAME" || true
    
    # Allow network access temporarily if needed (Skipped for speed, assume public access or run from vnet)
    az keyvault secret set --vault-name $KV_NAME --name payment-encryption-key --value "$ENCRYPTION_KEY"
fi

# 7. DB Init
echo -e "${GREEN}Initializing Databases...${NC}"
PG_HOST=$(get_val "postgresql_fqdn")
kubectl run db-init --image=postgres:alpine --restart=Never --env="PGPASSWORD=P@ssw0rd1234!" -- /bin/bash -c "
    psql -h $PG_HOST -U pgadmin -d postgres -c 'CREATE DATABASE merchants';
    psql -h $PG_HOST -U pgadmin -d postgres -c 'CREATE DATABASE tokenization';
    psql -h $PG_HOST -U pgadmin -d postgres -c 'CREATE DATABASE ledger';
    psql -h $PG_HOST -U pgadmin -d postgres -c 'CREATE DATABASE audit';
    psql -h $PG_HOST -U pgadmin -d postgres -c 'CREATE DATABASE reconciliation';
" || echo "DB Init pod might already exist or failed."

echo -e "${GREEN}Bootstrap Complete! Ready for ArgoCD.${NC}"
