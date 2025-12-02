provider "azurerm" {
  features {}
}

module "networking" {
  source              = "../../modules/networking"
  resource_group_name = var.resource_group_name
  location            = var.location
  environment         = var.environment
  hub_vnet_name       = var.hub_vnet_name
  spoke_vnet_name     = var.spoke_vnet_name
  hub_address_space   = var.hub_address_space
  spoke_address_space = var.spoke_address_space
  tags                = var.tags
}

module "aks" {
  source                    = "../../modules/aks"
  resource_group_name       = var.resource_group_name
  location                  = var.location
  aks_cluster_name          = "payments-aks-dev"
  aks_subnet_id             = module.networking.aks_subnet_id
  system_node_pool_vm_size  = "Standard_B2s"      # 2 vCPU, 4GB - BS Family
  user_node_pool_vm_size    = "Standard_D4s_v3"   # 4 vCPU, 16GB - DSv3 Family
  node_count_system         = 1
  node_count_user           = 1
  tags                      = var.tags
}

module "acr" {
  source              = "../../modules/acr"
  resource_group_name = var.resource_group_name
  location            = var.location
  acr_name            = "acrpaymentsdev${random_string.suffix.result}"  # Must be globally unique
  sku                 = "Basic"
  tags                = var.tags
}

# Get current Azure tenant ID
data "azurerm_client_config" "current" {}

module "keyvault" {
  source                           = "../../modules/keyvault"
  resource_group_name              = var.resource_group_name
  location                         = var.location
  keyvault_name                    = "kv-payments-${random_string.suffix.result}"  # Must be globally unique
  tenant_id                        = data.azurerm_client_config.current.tenant_id
  aks_kubelet_identity_object_id   = module.aks.kubelet_identity_object_id
  data_subnet_id                   = module.networking.data_subnet_id
  tags                             = var.tags
}

# Generate random suffix for globally unique names
resource "random_string" "suffix" {
  length  = 6
  special = false
  upper   = false
}

# Attach ACR to AKS
resource "azurerm_role_assignment" "aks_acr_pull" {
  scope                = module.acr.acr_id
  role_definition_name = "AcrPull"
  principal_id         = module.aks.kubelet_identity_object_id
}

# Private DNS Zone for PostgreSQL Flexible Server (standard zone name)
resource "azurerm_private_dns_zone" "postgres_dns" {
  name                = "privatelink.postgres.database.azure.com"
  resource_group_name = var.resource_group_name
  tags                = var.tags
}


resource "azurerm_private_dns_zone_virtual_network_link" "postgres_dns_link" {
  name                  = "link-postgres-spoke"
  resource_group_name   = var.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.postgres_dns.name
  virtual_network_id    = module.networking.spoke_vnet_id
  tags                  = var.tags
}

module "postgresql" {
  source                 = "../../modules/postgresql"
  resource_group_name    = var.resource_group_name
  location               = var.location
  postgresql_server_name = "psql-payments-${random_string.suffix.result}"
  data_subnet_id         = module.networking.postgres_subnet_id  # CORRECT
  private_dns_zone_id    = azurerm_private_dns_zone.postgres_dns.id
  admin_username         = "pgadmin"
  admin_password         = "P@ssw0rd1234!"
  tags                   = var.tags
}

module "redis" {
  source              = "../../modules/redis"
  resource_group_name = var.resource_group_name
  location            = var.location
  redis_name          = "redis-payments-${random_string.suffix.result}"
  data_subnet_id      = module.networking.data_subnet_id
  tags                = var.tags
}
