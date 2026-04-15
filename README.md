# Payments Sandbox Platform on Azure AKS

A production-grade, PCI-inspired B2B payment infrastructure on Azure Kubernetes Service (AKS).
**Core Philosophy:** "Stripe Lite" — Secure, Event-Driven, and Resilient.

## Project Status
**Current Phase:** ✅ Phase 8 (Security Hardening) COMPLETE
**Next Phase:** Phase 9 (Advanced Compliance & Failover) IN-PLANNING

### Recent Achievements
*   **Security Hardening (Phase 8):** Implemented micro-segmentation via `NetworkPolicy`, enforced `readOnlyRootFilesystem`, and integrated KMS using Azure Key Vault with Workload Identity.
*   **Zero-Trust Mesh:** Enforced **Strict mTLS** across the `payments-system` namespace using Istio PeerAuthentication.
*   **Observability Finalized:** End-to-end distributed tracing using OpenTelemetry + B3 propagation. All dashboards (Kiali, Grafana, Jaeger) accessible via stable Istio Ingress.
*   **ArgoCD Stability:** Resolved ingress redirect loops; GitOps-driven deployments are now fully operational.

## Key Features
*   **Zero-Trust Networking:** Strict isolation between simulation and production namespaces. Default-deny network policies with explicit whitelisting.
*   **Resilient Architecture:** 100% Go-based services with integrated circuit breakers, retries, and Redis-backed idempotency.
*   **KMS-Backed Tokenization:** Sensitive data is tokenized using AES-256 keys managed by Azure Key Vault.
*   **Distributed Tracing:** Full request waterfall visibility from API Gateway -> Payment Service -> Acquirer Simulator.
*   **Automated Bootstrap:** "One-Shot" platform setup via `scripts/bootstrap-cluster.sh`.

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
├── k8s-manifests/      # Kubernetes YAMLs (Base/Overlays pattern)
├── istio/              # Service Mesh Configs (Gateways, Security, VS)
├── services/           # Microservices (Go)
├── pkg/                # Shared Go Libraries (Resilience, Crypto, Tracing)
├── terraform/          # Infrastructure as Code (Azure)
├── scripts/            # Automation & Bootstrap Tools
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
    Run the automated script to configure namespaces, identities, and secrets.
    ```bash
    ./scripts/bootstrap-cluster.sh
    ```

3.  **Deploy via GitOps:**
    ```bash
    kubectl apply -f k8s-manifests/argocd-infra.yaml
    kubectl apply -f k8s-manifests/argocd-app.yaml
    ./scripts/configure-ingress.sh
    ```

## Future Roadmap (Phase 9)
*   **Secret Rotation:** Automated AES-256 key rotation in Key Vault via KEDA or Azure Functions.
*   **Compliance Auditing:** Enhanced `audit-service` to generate PCI-DSS ready reports.
*   **Global Failover:** Multi-region GSLB using Azure Traffic Manager.

## Tech Stack
- **Cloud:** Azure (AKS, Event Hubs, Key Vault, PostgreSQL)
- **Service Mesh:** Istio (mTLS Strict)
- **GitOps:** Argo CD
- **Observability:** OpenTelemetry, Prometheus, Grafana, Jaeger, Kiali
- **Languages:** Go (1.21+), React/TypeScript

## License
MIT (for portfolio/learning purposes)
