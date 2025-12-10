output "api_gateway_client_id" {
  description = "API Gateway Managed Identity Client ID"
  value       = azurerm_user_assigned_identity.api_gateway.client_id
}

output "payment_service_client_id" {
  description = "Payment Service Managed Identity Client ID"
  value       = azurerm_user_assigned_identity.payment_service.client_id
}

output "merchant_service_client_id" {
  description = "Merchant Service Managed Identity Client ID"
  value       = azurerm_user_assigned_identity.merchant_service.client_id
}

output "tokenization_service_client_id" {
  description = "Tokenization Service Managed Identity Client ID"
  value       = azurerm_user_assigned_identity.tokenization_service.client_id
}

output "acquirer_simulator_client_id" {
  description = "Acquirer Simulator Managed Identity Client ID"
  value       = azurerm_user_assigned_identity.acquirer_simulator.client_id
}

output "ledger_service_client_id" {
  description = "Ledger Service Managed Identity Client ID"
  value       = azurerm_user_assigned_identity.ledger_service.client_id
}

output "audit_service_client_id" {
  description = "Audit Service Managed Identity Client ID"
  value       = azurerm_user_assigned_identity.audit_service.client_id
}

output "reconciliation_service_client_id" {
  description = "Reconciliation Service Managed Identity Client ID"
  value       = azurerm_user_assigned_identity.reconciliation_service.client_id
}
