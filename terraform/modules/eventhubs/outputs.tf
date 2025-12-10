output "namespace_id" {
  description = "Event Hubs Namespace ID"
  value       = azurerm_eventhub_namespace.payments.id
}

output "namespace_name" {
  description = "Event Hubs Namespace name"
  value       = azurerm_eventhub_namespace.payments.name
}

output "namespace_fqdn" {
  description = "Event Hubs Namespace FQDN"
  value       = "${azurerm_eventhub_namespace.payments.name}.servicebus.windows.net"
}

output "payment_events_name" {
  description = "Payment events Event Hub name"
  value       = azurerm_eventhub.payment_events.name
}

output "reconciliation_events_name" {
  description = "Reconciliation events Event Hub name"
  value       = azurerm_eventhub.reconciliation_events.name
}
