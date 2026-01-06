# Payments Sandbox Platform Context

This `GEMINI.md` provides context for AI agents working on the "Payments Sandbox" project. It outlines the architecture, development workflow, and current status to ensure consistent and safe interactions.

## 1. Project Overview

**Goal:** Build a production-grade, PCI-inspired B2B payment infrastructure on Azure AKS.
**Core Philosophy:** "Stripe Lite" — a developer-focused payment gateway and ledger.
**Status:** Phase 5 (Resilience/GitOps) In Progress. Security Audit Complete: Implemented API Key Auth, Idempotency, and Secret Hardening.

## 2. Technical Architecture

### Stack
*   **Language:** Go (Fiber framework) for all backend services.
*   **Frontend:** React (Admin Dashboard - Planned).
*   **Database:** PostgreSQL (Azure Database for PostgreSQL in prod, Docker containers locally).
*   **Infrastructure:** Terraform (Azure VNet, AKS, Key Vault, Event Hubs).
*   **Orchestration:** Kubernetes (AKS), Istio Service Mesh (Strict mTLS).
*   **Security:** Network Policies (Zero Trust), Azure Workload Identity.
*   **Observability:** Prometheus, Grafana, Kiali, Jaeger.

### Microservices Map
| Service | Port | Description | Status |
| :--- | :--- | :--- | :--- |
| **API Gateway** | 8000 | Entry point, routing, auth. Go Fiber. | ✅ Running (Istio Ingress) |
| **Payment Service** | 8081 | Orchestrates payment flow (Auth/Capture). | ✅ Running (Strict mTLS) |
| **Tokenization** | 3003 | Handles sensitive PAN data. PCI scope boundary. | ✅ Running (Strict mTLS) |
| **Ledger Service** | 3005 | Double-entry bookkeeping. | ✅ Running (Redis Consumer) |
| **Audit Service** | 3006 | Immutable logs (HMAC signed). | ✅ Running (Redis Consumer) |
| **Acquirer Sim** | 3004 | Simulates bank responses. | ✅ Running |
| **Reconciliation**| 3007 | Batch settlement verification. | ✅ Running |

## 3. Master Recovery Plan (Cold Start Guide)

**CRITICAL:** Follow these steps precisely to restore the full environment after a `terraform destroy`.

### Phase 1: Infrastructure Provisioning
1.  **Deploy:** Run `terraform apply` in `terraform/envs/dev`.
2.  **Outputs:** Save outputs to `infrastructure-outputs.json`.
    ```bash
    terraform output -json > ../../../infrastructure-outputs.json
    ```

### Phase 2: Configuration & Build
1.  **Update Manifests:**
    *   Edit `k8s-manifests/overlays/dev-aks/kustomization.yaml`: Update `images` with the new ACR name.
    *   Edit `k8s-manifests/overlays/dev-aks/deployment-patch.yaml`: Update **ALL** `AZURE_CLIENT_ID` values and the `image` registry reference.
    *   Edit `k8s-manifests/overlays/dev-aks/configmap-patch.yaml`: Update `DB_HOST`, `REDIS_ADDR`, `KEYVAULT_URI` with new values.
2.  **Build & Push:**
    ```bash
    az aks get-credentials ...
    ./build-and-push.sh <NEW_ACR_NAME>
    ```

### Phase 3: Cluster Bootstrapping (The "Manual" Glue)
*Use `az aks command invoke` for all kubectl commands.*

1.  **Initialize Databases:**
    *   Terraform ONLY creates the `payments` DB. You must create the others.
    *   Command: `kubectl run db-init --image=postgres:alpine ...` (Use the script pattern to create `ledger`, `audit`, `tokenization`, `merchants`, `reconciliation`).
2.  **Seed Kubernetes Secrets:**
    *   Create `postgresql-credentials` and `redis-credentials` in `payments-system` namespace.
    *   Values: Fetch `redis_primary_access_key` from Azure CLI or Terraform state.

### Phase 4: Platform Layer (Istio & Observability)
1.  **Install Istio:**
    *   Use `helm template` locally to generate `istio-full-install.yaml` (base + istiod + ingress).
    *   Apply to cluster.
2.  **Install Observability:**
    *   Download standard addons (Prometheus, Kiali, Grafana, Jaeger) from Istio repo.
    *   Apply to `istio-system`.

### Phase 5: Application Deployment
1.  **Deploy:**
    ```bash
    kubectl kustomize k8s-manifests/overlays/dev-aks > full-deployment.yaml
    az aks command invoke ... --command "kubectl apply -f full-deployment.yaml"
    ```
2.  **Enable Mesh:**
    *   Label namespace: `kubectl label namespace payments-system istio-injection=enabled`.
    *   Restart deployments to inject sidecars.

### Phase 6: Security Hardening
1.  **mTLS:** Apply `istio/security/peer-authentication.yaml` and `istio/destinationrules/default-tls.yaml`.
2.  **Network Policies:** Apply `policies/networkpolicies/payments-system-policy.yaml`.

### Phase 7: Verification
1.  **Load Test:** Run the traffic generator loop (curl) inside a test pod in `payments-system`.
2.  **Verify:** Check `payment-service` logs for "Event published" and `ledger-service` API for balance updates.

## 4. Key Directives for AI Agents

1.  **Go Standardization:** All backend code MUST be Go 1.24+ using Fiber v2.
2.  **Security First:** Use `azidentity` for Azure Auth. **SSL (require)** is mandatory for DB connections in Azure.
3.  **Monorepo Build:** Docker builds MUST run from the root context.
4.  **Private Cluster:** Do not attempt direct `kubectl` calls; use `az aks command invoke`.

## 5. Directory Structure
```
/
├── apps/               # Frontend applications
├── ci-cd/              # CI/CD pipelines
├── k8s-manifests/      # Kubernetes YAMLs
│   ├── base/           # Base resources
│   └── overlays/       # Env-specific patches (CRITICAL: Update these on reset)
├── istio/              # Service Mesh Configs
│   ├── security/       # PeerAuthentication
│   ├── gateways/       # Ingress/Egress Gateways
│   └── destinationrules/ # mTLS settings
├── pkg/                # Shared Go Libraries
├── policies/           # Security Policies
│   └── networkpolicies/ # K8s NetworkPolicies (Zero Trust)
├── services/           # Microservices
├── terraform/          # Infrastructure as Code
└── build-and-push.sh   # Build script
```

## 6. Security Architecture Changes (Post-Audit)
*   **Authentication:** API Gateway now enforces `X-API-Key` validation by calling `merchant-service` internal endpoint.
*   **API Keys:** Switched from bcrypt to **SHA256** for performance/lookup efficiency.
*   **Idempotency:** `payment-service` uses Redis to enforce exactly-once processing via `Idempotency-Key` header.
*   **Secrets:** All service configs now **panic** if critical secrets (DB passwords, Keys) are missing. Default values like `postgres123` have been purged from code.
*   **Crypto:** Hardcoded fallback keys in `pkg/crypto` have been removed. Services require valid Azure Key Vault URI or explicit configuration.