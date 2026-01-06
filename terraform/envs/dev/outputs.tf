# Networking Outputs
output "spoke_vnet_id" {
  description = "Spoke VNet ID"
  value       = module.networking.spoke_vnet_id
}

output "aks_subnet_id" {
  description = "AKS Subnet ID"
  value       = module.networking.aks_subnet_id
}

# AKS Outputs
output "aks_cluster_name" {
  description = "AKS Cluster Name"
  value       = module.aks.aks_cluster_name
}

output "aks_oidc_issuer_url" {
  description = "AKS OIDC Issuer URL"
  value       = module.aks.oidc_issuer_url
}

# ACR Outputs
output "acr_login_server" {
  description = "ACR Login Server"
  value       = module.acr.acr_login_server
}

# Key Vault Outputs
output "keyvault_uri" {
  description = "Key Vault URI"
  value       = module.keyvault.keyvault_uri
}

# PostgreSQL Outputs
output "postgresql_fqdn" {
  description = "PostgreSQL Server FQDN"
  value       = module.postgresql.postgresql_fqdn
}

output "postgresql_server_id" {
  description = "PostgreSQL Server ID"
  value       = module.postgresql.postgresql_server_id
}

output "postgresql_database_name" {
  description = "PostgreSQL Database Name"
  value       = module.postgresql.postgresql_database_name
}

# Redis Outputs
output "redis_hostname" {
  description = "Redis Hostname"
  value       = module.redis.redis_hostname
}

output "redis_primary_access_key" {
  description = "Redis Primary Access Key"
  value       = module.redis.redis_primary_access_key
  sensitive   = true
}

# Event Hubs Outputs (NEW)
output "eventhubs_namespace_fqdn" {
  description = "Event Hubs Namespace FQDN"
  value       = module.eventhubs.namespace_fqdn
}

output "eventhubs_payment_events_name" {
  description = "Payment Events Event Hub Name"
  value       = module.eventhubs.payment_events_name
}

# Workload Identity Outputs (NEW)
output "api_gateway_client_id" {
  description = "API Gateway Managed Identity Client ID"
  value       = module.workload_identity.api_gateway_client_id
}

output "payment_service_client_id" {
  description = "Payment Service Managed Identity Client ID"
  value       = module.workload_identity.payment_service_client_id
}

output "merchant_service_client_id" {
  description = "Merchant Service Managed Identity Client ID"
  value       = module.workload_identity.merchant_service_client_id
}

output "tokenization_service_client_id" {
  description = "Tokenization Service Managed Identity Client ID"
  value       = module.workload_identity.tokenization_service_client_id
}

output "acquirer_simulator_client_id" {
  description = "Acquirer Simulator Managed Identity Client ID"
  value       = module.workload_identity.acquirer_simulator_client_id
}

output "ledger_service_client_id" {
  description = "Ledger Service Managed Identity Client ID"
  value       = module.workload_identity.ledger_service_client_id
}

output "audit_service_client_id" {
  description = "Audit Service Managed Identity Client ID"
  value       = module.workload_identity.audit_service_client_id
}

output "reconciliation_service_client_id" {
  description = "Reconciliation Service Managed Identity Client ID"
  value       = module.workload_identity.reconciliation_service_client_id
}
