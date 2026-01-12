# Payments Sandbox Platform Context

This `GEMINI.md` provides context for AI agents working on the "Payments Sandbox" project. It outlines the architecture, development workflow, and current status to ensure consistent and safe interactions.

## 1. Project Overview

**Goal:** Build a production-grade, PCI-inspired B2B payment infrastructure on Azure AKS.
**Core Philosophy:** "Stripe Lite" — a developer-focused payment gateway and ledger.
**Status:** Phase 6 (Resilience Policies) COMPLETE. Phase 7 (Observability & Verification) READY.
**Key Achievement:** Implemented Istio resilience policies (Retries, Timeouts, Circuit Breakers) for core services.

## 2. Technical Architecture (Post-Phase 5)

### Stack Updates
*   **GitOps:** ArgoCD manages the cluster via two streams: `platform-infra` and `payments-platform`.
*   **Service Mesh:** Istio (Strict mTLS) with standardized port naming (`http-api`) for L7 metrics/tracing.
*   **Identity:** Azure Workload Identity for all services.
*   **KMS:** `pkg/crypto` refactored to use Key Vault **Secrets** for Symmetric AES keys (compatible with Standard tier).

### Microservices Map
All services use port `http-api` for Istio compatibility.
| Service | Internal Port | Description |
| :--- | :--- | :--- |
| **API Gateway** | 3000 (Svc 80) | Entry point, routing, auth, circuit breakers. |
| **Payment Service** | 8081 | Orchestrates payment flow (Auth/Capture). |
| **Tokenization** | 3003 | Handles sensitive PAN data. Uses KV Secret. |
| **Ledger Service** | 3005 | Double-entry bookkeeping. |
| **Audit Service** | 3006 | Immutable logs (HMAC signed). |
| **Merchant Service**| 3002 | Merchant and API Key management. |
| **Acquirer Sim** | 3004 | Simulates bank responses. |
| **Reconciliation**| 3007 | Batch settlement verification. |

## 3. Master Recovery Plan (Cold Start Guide)

**CRITICAL:** Follow these steps to restore the environment after a `terraform destroy`.

### Phase 1: Infrastructure
1.  **Deploy:** Run `terraform apply -var-file="dev.tfvars"` in `terraform/envs/dev`.
2.  **Outputs:** Save outputs: `terraform output -json > ../../../infrastructure-outputs.json`.

### Phase 2: Cluster Bootstrapping (Manual Glue)
*Use `az aks command invoke` for these commands.*

1.  **Create Namespaces:** `argocd`, `payments-system`, `payments-data`, `istio-system`.
2.  **Initialize Databases:** Run a `postgres:alpine` pod to create: `merchants`, `tokenization`, `ledger`, `audit`, `reconciliation`.
3.  **Seed Kubernetes Secrets:** 
    *   `postgresql-credentials` (username, password, host, port, sslmode).
    *   `redis-credentials` (password).
    *   `audit-secrets` (AUDIT_HMAC_KEY).
4.  **Seed Key Vault Secret:** Create a Secret (not Key) named `payment-encryption-key` in Key Vault with a 32-byte base64 string.

### Phase 3: GitOps Activation
1.  **ArgoCD Install:** Apply the official install manifest to the `argocd` namespace.
2.  **Infrastructure App:** Apply `k8s-manifests/argocd-infra.yaml`. Wait for Istio pods to be `Running`.
3.  **Application App:** Apply `k8s-manifests/argocd-app.yaml`. This will deploy all 8 services.

## 4. Key Directives for AI Agents

1.  **Port Naming:** Service ports MUST be named `http-api` (or start with `http-`) or tracing/metrics will fail.
2.  **Workload Identity:** New deployments MUST include the `azure.workload.identity/use: "true"` label and reference a ServiceAccount with the `azure.workload.identity/client-id` annotation.
3.  **Monorepo Build:** Docker builds MUST run from the root context using `build-and-push.sh`.
4.  **Private Cluster:** Use `az aks command invoke` for all `kubectl` operations.

## 5. Security Architecture
*   **mTLS:** Enforced via `PeerAuthentication` (Strict).
*   **Vault Access:** Identity `id-tokenization-service-dev` requires `Key Vault Secrets User` role on the vault.
*   **Ingress:** Dashboards are exposed via `dashboard-gateway` using `nip.io` subdomains on the Public Ingress IP.
