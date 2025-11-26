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
  source              = "../../modules/aks"
  resource_group_name = var.resource_group_name
  location            = var.location
  aks_cluster_name    = "payments-aks-dev"
  aks_subnet_id       = module.networking.aks_subnet_id
  node_pool_vm_size   = "Standard_D2s_v3"
  node_count_system   = 1
  node_count_user     = 2
  tags                = var.tags
}
