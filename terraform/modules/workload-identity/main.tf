# User-Assigned Managed Identity for API Gateway
resource "azurerm_user_assigned_identity" "api_gateway" {
  name                = "id-api-gateway-${var.environment}"
  resource_group_name = var.resource_group_name
  location            = var.location
  tags                = var.tags
}

# Federated Credential for API Gateway
resource "azurerm_federated_identity_credential" "api_gateway" {
  name                = "api-gateway-federated"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.api_gateway.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = var.oidc_issuer_url
  subject             = "system:serviceaccount:payments-system:api-gateway"
}

# User-Assigned Managed Identity for Payment Service
resource "azurerm_user_assigned_identity" "payment_service" {
  name                = "id-payment-service-${var.environment}"
  resource_group_name = var.resource_group_name
  location            = var.location
  tags                = var.tags
}

resource "azurerm_federated_identity_credential" "payment_service" {
  name                = "payment-service-federated"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.payment_service.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = var.oidc_issuer_url
  subject             = "system:serviceaccount:payments-system:payment-service"
}

# User-Assigned Managed Identity for Merchant Service
resource "azurerm_user_assigned_identity" "merchant_service" {
  name                = "id-merchant-service-${var.environment}"
  resource_group_name = var.resource_group_name
  location            = var.location
  tags                = var.tags
}

resource "azurerm_federated_identity_credential" "merchant_service" {
  name                = "merchant-service-federated"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.merchant_service.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = var.oidc_issuer_url
  subject             = "system:serviceaccount:payments-system:merchant-service"
}

# User-Assigned Managed Identity for Tokenization Service
resource "azurerm_user_assigned_identity" "tokenization_service" {
  name                = "id-tokenization-service-${var.environment}"
  resource_group_name = var.resource_group_name
  location            = var.location
  tags                = var.tags
}

resource "azurerm_federated_identity_credential" "tokenization_service" {
  name                = "tokenization-service-federated"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.tokenization_service.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = var.oidc_issuer_url
  subject             = "system:serviceaccount:payments-system:tokenization-service"
}

# User-Assigned Managed Identity for Acquirer Simulator
resource "azurerm_user_assigned_identity" "acquirer_simulator" {
  name                = "id-acquirer-simulator-${var.environment}"
  resource_group_name = var.resource_group_name
  location            = var.location
  tags                = var.tags
}

resource "azurerm_federated_identity_credential" "acquirer_simulator" {
  name                = "acquirer-simulator-federated"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.acquirer_simulator.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = var.oidc_issuer_url
  subject             = "system:serviceaccount:payments-system:acquirer-simulator"
}

# User-Assigned Managed Identity for Ledger Service
resource "azurerm_user_assigned_identity" "ledger_service" {
  name                = "id-ledger-service-${var.environment}"
  resource_group_name = var.resource_group_name
  location            = var.location
  tags                = var.tags
}

resource "azurerm_federated_identity_credential" "ledger_service" {
  name                = "ledger-service-federated"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.ledger_service.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = var.oidc_issuer_url
  subject             = "system:serviceaccount:payments-data:ledger-service"
}

# User-Assigned Managed Identity for Audit Service
resource "azurerm_user_assigned_identity" "audit_service" {
  name                = "id-audit-service-${var.environment}"
  resource_group_name = var.resource_group_name
  location            = var.location
  tags                = var.tags
}

resource "azurerm_federated_identity_credential" "audit_service" {
  name                = "audit-service-federated"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.audit_service.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = var.oidc_issuer_url
  subject             = "system:serviceaccount:payments-data:audit-service"
}

# User-Assigned Managed Identity for Reconciliation Service
resource "azurerm_user_assigned_identity" "reconciliation_service" {
  name                = "id-reconciliation-service-${var.environment}"
  resource_group_name = var.resource_group_name
  location            = var.location
  tags                = var.tags
}

resource "azurerm_federated_identity_credential" "reconciliation_service" {
  name                = "reconciliation-service-federated"
  resource_group_name = var.resource_group_name
  parent_id           = azurerm_user_assigned_identity.reconciliation_service.id
  audience            = ["api://AzureADTokenExchange"]
  issuer              = var.oidc_issuer_url
  subject             = "system:serviceaccount:payments-data:reconciliation-service"
}

# RBAC: Event Hubs Data Sender (Payment Service)
resource "azurerm_role_assignment" "payment_eventhubs_sender" {
  scope                = var.eventhubs_namespace_id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azurerm_user_assigned_identity.payment_service.principal_id
}

# RBAC: Event Hubs Data Sender (Reconciliation Service)
resource "azurerm_role_assignment" "reconciliation_eventhubs_sender" {
  scope                = var.eventhubs_namespace_id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azurerm_user_assigned_identity.reconciliation_service.principal_id
}

# RBAC: Event Hubs Data Receiver (Ledger Service)
resource "azurerm_role_assignment" "ledger_eventhubs_receiver" {
  scope                = var.eventhubs_namespace_id
  role_definition_name = "Azure Event Hubs Data Receiver"
  principal_id         = azurerm_user_assigned_identity.ledger_service.principal_id
}

# RBAC: Event Hubs Data Receiver (Audit Service)
resource "azurerm_role_assignment" "audit_eventhubs_receiver" {
  scope                = var.eventhubs_namespace_id
  role_definition_name = "Azure Event Hubs Data Receiver"
  principal_id         = azurerm_user_assigned_identity.audit_service.principal_id
}

# RBAC: Key Vault Crypto User (Tokenization Service)
resource "azurerm_role_assignment" "tokenization_keyvault_crypto" {
  scope                = var.keyvault_id
  role_definition_name = "Key Vault Crypto User"
  principal_id         = azurerm_user_assigned_identity.tokenization_service.principal_id
}

# RBAC: Key Vault Secrets User (API Gateway - for JWT signing)
resource "azurerm_role_assignment" "api_gateway_keyvault_secrets" {
  scope                = var.keyvault_id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_user_assigned_identity.api_gateway.principal_id
}
