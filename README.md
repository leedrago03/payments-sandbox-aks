# Payments Sandbox Platform on Azure AKS

A production-grade, PCI-inspired B2B payment infrastructure on Azure Kubernetes Service (AKS).
**Core Philosophy:** "Stripe Lite" — Secure, Event-Driven, and Resilient.

## Project Status
**Current Phase:** Phase 7 (Verification & Observability) IN-PROGRESS
**Previous Phase:** Phase 6 (Resilience) COMPLETE

### Recent Achievements
*   **"One-Shot" Enterprise Bootstrap:** Fully automated cluster setup via `scripts/bootstrap-cluster.sh`.
*   **Resilience Verified:** System handles 100 concurrent requests and survives 50% upstream failure rates (Chaos Testing).
*   **Service Mesh:** All microservices running with Istio sidecars (mTLS Strict).

## Key Features
*   **Zero-Trust Networking:** Strict mTLS between all microservices via Istio.
*   **API Security:** SHA256-hashed API Key authentication enforced at the Gateway.
*   **Event-Driven:** Async processing using Azure Event Hubs.
*   **Idempotency:** Redis-backed idempotency keys for Payment Service reliability.
*   **Secret Management:** No hardcoded secrets. Azure Key Vault integration with Workload Identity.
*   **Observability:** Prometheus, Grafana, Kiali, and Jaeger (Integration in progress).

## Architecture
Based on [Azure AKS Secure Baseline](https://learn.microsoft.com/en-us/azure/architecture/reference-architectures/containers/aks/baseline-aks).

### Microservices Map
| Service | Port | Description |
| :--- | :--- | :--- |
| **API Gateway** | 8000 | Entry point, routing, auth. Go Fiber. |
| **Payment Service** | 8081 | Orchestrates payment flow (Auth/Capture). |
| **Tokenization** | 3003 | Handles sensitive PAN data. PCI scope boundary. |
| **Ledger Service** | 3005 | Double-entry bookkeeping. |
| **Audit Service** | 3006 | Immutable logs (HMAC signed). |
| **Merchant Service** | 3002 | Merchant profile & API key management. |
| **Acquirer Sim** | 3004 | Simulates bank responses. |
| **Reconciliation**| 3007 | Batch settlement verification. |

## Directory Structure
```
├── apps/               # Frontend applications
├── ci-cd/              # CI/CD pipelines
├── k8s-manifests/      # Kubernetes YAMLs (Base/Overlays pattern)
├── istio/              # Service Mesh Configs
├── pkg/                # Shared Go Libraries (Crypto, Events, Logging)
├── scripts/            # Automation & Bootstrap Tools
├── services/           # Microservices (Go)
├── terraform/          # Infrastructure as Code
└── GEMINI.md           # AI Context & Master Plan
```

## Quick Start (Enterprise Standard)

**Prerequisite:** Azure CLI, Terraform, Kubectl, Helm.

1.  **Provision Infrastructure:**
    ```bash
    cd terraform/envs/dev
    terraform apply -var-file="dev.tfvars"
    terraform output -json > ../../../infrastructure-outputs.json
    ```

2.  **Bootstrap Platform:**
    Run the "One-Shot" script to configure namespaces, secrets, DBs, and Key Vault.
    ```bash
    cd ../../../
    chmod +x scripts/bootstrap-cluster.sh
    ./scripts/bootstrap-cluster.sh
    ```

3.  **Deploy via GitOps:**
    ```bash
    kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
    kubectl apply -f k8s-manifests/argocd-infra.yaml
    kubectl apply -f k8s-manifests/argocd-app.yaml
    ```

## Known Issues (Phase 7)
*   **ArgoCD:** Redirect loop on ingress (Workaround: Use port-forward).
*   **Jaeger:** Traces missing for some services (Fix in progress).
*   **Kiali:** RBAC permissions need adjustment for full graph visibility.

## Tech Stack
- **Cloud:** Azure
- **Kubernetes:** AKS (private, workload identity)
- **Service Mesh:** Istio
- **Events:** Azure Event Hubs
- **Data:** Azure PostgreSQL, Redis, Azure Blob Storage
- **Secrets:** Azure Key Vault
- **GitOps:** Argo CD
- **Observability:** Prometheus, Grafana, Jaeger, Kiali
- **Languages:** Go (1.21+), React/TypeScript

## License
MIT (for portfolio/learning purposes)