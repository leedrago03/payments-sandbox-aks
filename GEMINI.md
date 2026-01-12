# Payments Sandbox Platform Context

This `GEMINI.md` provides the authoritative context for AI agents working on the "Payments Sandbox" project. Adhere strictly to the architectural standards and recovery plans outlined here.

## 1. Project Overview
**Goal:** Build a production-grade, PCI-inspired B2B payment infrastructure on Azure AKS.
**Core Philosophy:** "Stripe Lite" — Secure, Event-Driven, and Resilient.
**Status:** Phase 6 (Resilience) COMPLETE. Phase 7 (Verification & Observability) IN-PROGRESS.
**Key Achievement:** Implemented "One-Shot" Enterprise Bootstrap; verified mesh resilience via stress and chaos testing.

## 2. Technical Architecture (Post-Phase 7 Refactor)

### Standardized Manifest Structure
To prevent configuration drift and "patching hell," the following rules are enforced:
*   **Base Manifests (`k8s-manifests/base`):** MUST use placeholder images (e.g., `image: payment-service`) and generic service names. No environment-specific URLs or Client IDs.
*   **Overlay Manifests (`k8s-manifests/overlays/dev-aks`):** Use `kustomization.yaml` to replace images.
*   **Patches (`deployment-patch.yaml`):** MUST NOT contain `image:` fields. Use only for environment-specific metadata like `AZURE_CLIENT_ID` or specific labels.

### Automation Tooling
*   `scripts/configure-manifests.py`: Automatically syncs `infrastructure-outputs.json` (from Terraform) into K8s manifests (URIs, Client IDs).
*   `scripts/bootstrap-cluster.sh`: The "Enterprise One-Shot" script. Sets namespaces, seeds secrets, initializes 5 logical DBs, and configures Key Vault.
*   `scripts/configure-ingress.sh`: Dynamically detects LoadBalancer IP and updates `nip.io` ingress routes for dashboards.

---

## 3. Master Recovery Plan (Enterprise Standard)

**CRITICAL:** Follow these steps to restore the environment after a `terraform destroy`.

### Phase 1: Infrastructure
1.  **Provision:** `cd terraform/envs/dev && terraform apply -var-file="dev.tfvars"`.
2.  **Capture:** `terraform output -json > ../../../infrastructure-outputs.json`.

### Phase 2: Platform Bootstrapping
Run the automated bootstrap from the root directory:
```bash
chmod +x scripts/bootstrap-cluster.sh
./scripts/bootstrap-cluster.sh
```

### Phase 3: GitOps & Deployment
1.  **Install ArgoCD:** `kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml`.
2.  **App Infrastructure:** `kubectl apply -f k8s-manifests/argocd-infra.yaml`.
3.  **App Services:** `kubectl apply -f k8s-manifests/argocd-app.yaml`.
4.  **Configure Ingress:** `chmod +x scripts/configure-ingress.sh && ./scripts/configure-ingress.sh`.

---

## 4. Current State & Known Issues (Phase 7)

### Verified Capabilities
*   **Resilience:** System verified to handle 100 concurrent requests.
*   **Chaos Handling:** Fault injection (50% error rate on Acquirer) confirmed that `payment-service` handles upstream failures gracefully without cascading crashes.
*   **Connectivity:** All 8 services are Running (2/2) with Istio sidecars. Public access via Ingress Gateway (Port 80/443) is enabled.

### Open Issues (Continue from here)
1.  **ArgoCD (`ERR_TOO_MANY_REDIRECTS`):** Despite `--insecure` flag and `X-Forwarded-Proto` header injection, the redirect loop between Gateway (SSL) and ArgoCD (Plain) persists. Needs deeper Envoy filter or ArgoCD config map audit.
2.  **Jaeger (Empty Traces):** Zipkin enabled in MeshConfig and workloads restarted, but application spans (API Gateway -> Payment) are still missing. Likely requires `telemetry` CRD configuration or header propagation verification in Go code.
3.  **Kiali Permissions:** Occasional "Forbidden" errors when fetching deployment status. Needs ClusterRole check.

## 5. Security Architecture
*   **mTLS:** Enforced (Strict) across `payments-system`.
*   **Workload Identity:** Azure identities mapped to K8s ServiceAccounts.
*   **KMS:** Tokenization service uses AES-256 Symmetric keys stored as **Secrets** in Azure Key Vault (Standard Tier).
*   **Networking:** Port 443 (HTTPS), 80 (HTTP), and 3000 (Gateway) allowed via `nsg-aks-dev`.