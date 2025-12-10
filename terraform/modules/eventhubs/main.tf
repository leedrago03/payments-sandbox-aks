# Azure Event Hubs Namespace
resource "azurerm_eventhub_namespace" "payments" {
  name                = "ehns-payments-${var.environment}"
  resource_group_name = var.resource_group_name
  location            = var.location
  sku                 = "Standard"
  capacity            = 1
  
  # Network isolation
  public_network_access_enabled = false
  
  tags = var.tags
}

# Event Hub: Payment Events
resource "azurerm_eventhub" "payment_events" {
  name                = "payment-events"
  namespace_name      = azurerm_eventhub_namespace.payments.name
  resource_group_name = var.resource_group_name
  partition_count     = 2
  message_retention   = 1  # 1 day retention
}

# Consumer Group for Ledger Service
resource "azurerm_eventhub_consumer_group" "ledger_consumer" {
  name                = "ledger-consumer"
  namespace_name      = azurerm_eventhub_namespace.payments.name
  eventhub_name       = azurerm_eventhub.payment_events.name
  resource_group_name = var.resource_group_name
}

# Consumer Group for Audit Service
resource "azurerm_eventhub_consumer_group" "audit_consumer" {
  name                = "audit-consumer"
  namespace_name      = azurerm_eventhub_namespace.payments.name
  eventhub_name       = azurerm_eventhub.payment_events.name
  resource_group_name = var.resource_group_name
}

# Event Hub: Reconciliation Events
resource "azurerm_eventhub" "reconciliation_events" {
  name                = "reconciliation-events"
  namespace_name      = azurerm_eventhub_namespace.payments.name
  resource_group_name = var.resource_group_name
  partition_count     = 1
  message_retention   = 1
}

# Consumer Group for Alert Service
resource "azurerm_eventhub_consumer_group" "alerts_consumer" {
  name                = "alerts-consumer"
  namespace_name      = azurerm_eventhub_namespace.payments.name
  eventhub_name       = azurerm_eventhub.reconciliation_events.name
  resource_group_name = var.resource_group_name
}

# Private Endpoint for Event Hubs
resource "azurerm_private_endpoint" "eventhubs_pe" {
  name                = "pe-eventhubs-${var.environment}"
  location            = var.location
  resource_group_name = var.resource_group_name
  subnet_id           = var.data_subnet_id

  private_service_connection {
    name                           = "psc-eventhubs-${var.environment}"
    private_connection_resource_id = azurerm_eventhub_namespace.payments.id
    is_manual_connection           = false
    subresource_names              = ["namespace"]
  }

  tags = var.tags
}

# Private DNS Zone for Event Hubs
resource "azurerm_private_dns_zone" "eventhubs_dns" {
  name                = "privatelink.servicebus.windows.net"
  resource_group_name = var.resource_group_name
  tags                = var.tags
}

# Link DNS Zone to VNet
resource "azurerm_private_dns_zone_virtual_network_link" "eventhubs_dns_link" {
  name                  = "link-eventhubs-${var.environment}"
  resource_group_name   = var.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.eventhubs_dns.name
  virtual_network_id    = var.vnet_id
  tags                  = var.tags
}

# DNS A Record for Private Endpoint
resource "azurerm_private_dns_a_record" "eventhubs_dns_record" {
  name                = azurerm_eventhub_namespace.payments.name
  zone_name           = azurerm_private_dns_zone.eventhubs_dns.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [azurerm_private_endpoint.eventhubs_pe.private_service_connection[0].private_ip_address]
  tags                = var.tags
}
