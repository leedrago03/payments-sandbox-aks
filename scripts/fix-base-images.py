import os
import yaml

base_dir = 'k8s-manifests/base'
services = [
    'payment-service', 'tokenization-service', 'audit-service', 
    'ledger-service', 'merchant-service', 'acquirer-simulator', 
    'reconciliation-service', 'api-gateway'
]

for service in services:
    path = os.path.join(base_dir, service)
    # Find deployment file
    if os.path.exists(os.path.join(path, 'deployment.yaml')):
        deploy_file = os.path.join(path, 'deployment.yaml')
    elif os.path.exists(os.path.join(path, 'resources.yaml')):
        deploy_file = os.path.join(path, 'resources.yaml')
    else:
        print(f"No deployment file found for {service}")
        continue
        
    print(f"Fixing {deploy_file}...")
    
    # Read all docs
    with open(deploy_file, 'r') as f:
        docs = list(yaml.safe_load_all(f))
        
    for doc in docs:
        if doc.get('kind') == 'Deployment':
            # Check containers
            for container in doc['spec']['template']['spec']['containers']:
                # The container name usually matches the service name
                # We want to set image = service name to allow Kustomize replacement
                container['image'] = service 
                
    with open(deploy_file, 'w') as f:
        yaml.dump_all(docs, f, default_flow_style=False, sort_keys=False)

print("Base manifests updated to use placeholder images.")
