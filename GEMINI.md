# Payments Sandbox Platform Context

This `GEMINI.md` provides the authoritative context for AI agents working on the "Payments Sandbox" project. Adhere strictly to the architectural standards and recovery plans outlined here.

## 1. Project Overview
**Goal:** Build a production-grade, PCI-inspired B2B payment infrastructure on Azure AKS.
**Core Philosophy:** "Stripe Lite" — Secure, Event-Driven, and Resilient.
**Status:** 
*   ✅ Phase 6 (Resilience) COMPLETE.
*   ✅ Phase 7 (Verification & Observability) COMPLETE.
*   ✅ Phase 8 (Security Hardening) COMPLETE.

## 2. Technical Architecture

### Observability Stack (Phase 7 Finalized)
*   **Distributed Tracing:** Implemented via OpenTelemetry + B3 Propagation. All Go services propagate `X-B3` headers.
*   **Istio Mesh:** Configured with `Telemetry` CRD for 100% Zipkin-style sampling.
*   **Dashboards:** Kiali, Grafana, and Jaeger accessible via Istio Ingress with fixed RBAC and SSL redirect logic.

### Security Architecture (Phase 8 Hardened)
*   **Micro-Segmentation:** Fine-grained `NetworkPolicy` per service. Default-deny with explicit whitelisting of neighbors (e.g., Gateway -> Payment only).
*   **Pod Hardening:** All deployments enforce `readOnlyRootFilesystem: true` with `emptyDir` mounts for `/tmp`.
*   **mTLS:** Enforced (Strict) across `payments-system`.
*   **Workload Identity:** Azure identities mapped to K8s ServiceAccounts.
*   **KMS:** Tokenization service uses AES-256 Symmetric keys stored as **Secrets** in Azure Key Vault.

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

## 4. Current State & Verified Capabilities

### Verified Capabilities
*   **Zero-Trust Mesh:** Strict isolation between simulation and production namespaces.
*   **End-to-End Tracing:** Complete request waterfall from Gateway -> Payment -> Acquirer.
*   **ArgoCD:** Stable ingress without redirect loops.
*   **Resilience:** System verified to handle 100 concurrent requests with graceful degradation.

### Open Items (Phase 9 - Future)
1.  **Secret Rotation:** Implement automated rotation for the AES-256 keys in Key Vault via Azure Functions or KEDA.
2.  **Compliance Audit Logs:** Enhance `audit-service` to generate PCI-ready reports.
3.  **Cross-Region Failover:** Explore Azure Traffic Manager for multi-region GSLB.
