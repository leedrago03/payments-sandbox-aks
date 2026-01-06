# Payments Sandbox Platform on Azure AKS

A production-like, PCI-inspired payment processing sandbox built on Azure Kubernetes Service (AKS) with enterprise security, observability, and GitOps.

## Project Status
✅ **Phase 5 In Progress** - Resilience, GitOps & Security Hardening
(Security Audit Complete: Jan 2026)

## Key Security Features (Hardened)
- **Zero-Trust Networking:** Strict mTLS between all microservices via Istio.
- **API Security:** SHA256-hashed API Key authentication enforced at the Gateway.
- **Idempotency:** Redis-backed idempotency keys for Payment Service reliability.
- **Secret Management:** No hardcoded secrets. All services require Azure Key Vault injection or explicit env vars; app panics on missing configuration.
- **Least Privilege:** Workload Identity for all Azure interactions.

## Overview
This project demonstrates end-to-end platform engineering skills by building a realistic payment gateway and ledger system with:
- Private AKS cluster in hub-spoke network topology
- Istio service mesh with mTLS
- Tokenization, payments, double-entry ledger, reconciliation
- Azure-native integration (Key Vault, Event Hubs, PostgreSQL)
- Hardened CI/CD with image signing and policy enforcement
- Full observability stack (Prometheus, Grafana, Jaeger, Loki)

## Architecture
Based on [Azure AKS Secure Baseline](https://learn.microsoft.com/en-us/azure/architecture/reference-architectures/containers/aks/baseline-aks) and [AKS PCI-DSS guidance](https://learn.microsoft.com/en-us/azure/aks/pci-intro).

### Microservices Map
| Service | Port | Description |
| :--- | :--- | :--- |
| **API Gateway** | 8000 | Entry point, routing, auth. Go Fiber. |
| **Payment Service** | 8081 | Orchestrates payment flow (Auth/Capture). |
| **Tokenization** | 3003 | Handles sensitive PAN data. PCI scope boundary. |
| **Ledger Service** | 3005 | Double-entry bookkeeping. |
| **Audit Service** | 3006 | Immutable logs (HMAC signed). |
| **Acquirer Sim** | 3004 | Simulates bank responses. |
| **Reconciliation**| 3007 | Batch settlement verification. |

## Repository Structure
├── apps/# Microservices (Go + React)

├── terraform/# Infrastructure as Code

├── k8s-manifests/ # Kubernetes manifests (GitOps source)

├── istio/ # Service mesh configuration

├── policies/ # OPA Gatekeeper & Network Policies

├── ci-cd/ # CI/CD pipelines

└── docs/ # Architecture & runbooks

## Quick Start
See [docs/00-readme/local-setup.md](docs/00-readme/local-setup.md) for prerequisites and setup instructions.

## Documentation
- [Architecture Overview](docs/architecture/phase0-architecture-notes.md)
- [Environments](docs/architecture/environments.md)
- [Terraform Design](docs/architecture/phase1-terraform-design.md)

## Tech Stack
- **Cloud:** Azure
- **Kubernetes:** AKS (private, workload identity)
- **Service Mesh:** Istio
- **Events:** Azure Event Hubs
- **Data:** Azure PostgreSQL, Azure Blob Storage
- **Secrets:** Azure Key Vault
- **GitOps:** Argo CD
- **Observability:** Prometheus, Grafana, Jaeger, Loki
- **Policy:** OPA Gatekeeper
- **Languages:** Go, React/TypeScript


## Development Access (Optional)

For easier cluster management, you can deploy a jumpbox VM:

**When to use:**
- You need frequent kubectl/helm access
- Debugging and testing microservices
- Learning/development environments

**When to skip:**
- Cost optimization for students
- You have existing VPN access
- Using CI/CD pipelines only

**To deploy:**
cd terraform/envs/dev-tools
terraform apply


## License
MIT (for portfolio/learning purposes)
