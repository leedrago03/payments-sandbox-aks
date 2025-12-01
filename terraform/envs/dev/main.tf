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
