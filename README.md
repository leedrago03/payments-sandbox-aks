# Payments Sandbox Platform on Azure AKS

A production-like, PCI-inspired payment processing sandbox built on Azure Kubernetes Service (AKS) with enterprise security, observability, and GitOps.

## Project Status
🚧 **Phase 1 In Progress** - Azure & AKS Foundation (IaC)

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

## Repository Structure
├── apps/ # Microservices (Go + React)
├── terraform/ # Infrastructure as Code
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

## License
MIT (for portfolio/learning purposes)
