import json
import yaml
import os

# Load outputs
with open('infrastructure-outputs.json', 'r') as f:
    outputs = json.load(f)

def get_val(key):
    return outputs[key]['value']

acr_server = get_val('acr_login_server')
kv_uri = get_val('keyvault_uri')
db_host = get_val('postgresql_fqdn')
redis_host = get_val('redis_hostname') + ":6380"

# Client IDs
ids = {
    'api-gateway': get_val('api_gateway_client_id'),
    'payment-service': get_val('payment_service_client_id'),
    'tokenization-service': get_val('tokenization_service_client_id'),
    'ledger-service': get_val('ledger_service_client_id'),
    'audit-service': get_val('audit_service_client_id'),
    'merchant-service': get_val('merchant_service_client_id'),
    'acquirer-simulator': get_val('acquirer_simulator_client_id'),
    'reconciliation-service': get_val('reconciliation_service_client_id')
}

# 1. Update kustomization.yaml (Images)
kust_path = 'k8s-manifests/overlays/dev-aks/kustomization.yaml'
with open(kust_path, 'r') as f:
    kust = yaml.safe_load(f)

for img in kust['images']:
    # Format: acrserver/payments-sandbox/service-name
    # We replace the registry part
    service_name = img['name']
    img['newName'] = f"{acr_server}/payments-sandbox/{service_name}"

with open(kust_path, 'w') as f:
    yaml.dump(kust, f, default_flow_style=False, sort_keys=False)

# 2. Update configmap-patch.yaml
cm_path = 'k8s-manifests/overlays/dev-aks/configmap-patch.yaml'
with open(cm_path, 'r') as f:
    cms = list(yaml.safe_load_all(f))

for cm in cms:
    if 'data' in cm:
        if 'DB_HOST' in cm['data']:
            cm['data']['DB_HOST'] = db_host
        if 'REDIS_ADDR' in cm['data']:
            cm['data']['REDIS_ADDR'] = redis_host
        if 'KEYVAULT_URI' in cm['data']:
            cm['data']['KEYVAULT_URI'] = kv_uri
        # ACR_SERVER might not be in configmaps usually, but if needed:
        if 'ACR_SERVER' in cm['data']:
             cm['data']['ACR_SERVER'] = acr_server

with open(cm_path, 'w') as f:
    yaml.dump_all(cms, f, default_flow_style=False, sort_keys=False)

# 3. Update deployment-patch.yaml (Env Vars + Image Registry for some if hardcoded?)
# The deployment patch primarily sets AZURE_CLIENT_ID env vars
deploy_path = 'k8s-manifests/overlays/dev-aks/deployment-patch.yaml'
with open(deploy_path, 'r') as f:
    patches = list(yaml.safe_load_all(f))

new_patches = []
for patch in patches:
    if patch['kind'] == 'Deployment':
        name = patch['metadata']['name']
        # Find the container spec
        containers = patch['spec']['template']['spec']['containers']
        for container in containers:
            if 'env' in container:
                for env in container['env']:
                    if env['name'] == 'AZURE_CLIENT_ID':
                        if name in ids:
                            env['value'] = ids[name]
                    elif env['name'] == 'AZURE_KEY_VAULT_URI':
                        env['value'] = kv_uri
        
    new_patches.append(patch)

with open(deploy_path, 'w') as f:
    yaml.dump_all(new_patches, f, default_flow_style=False, sort_keys=False)

print("Manifests updated successfully.")
