resource "azurerm_kubernetes_cluster" "aks" {
  name                = var.aks_cluster_name
  location            = var.location
  resource_group_name = var.resource_group_name
  dns_prefix          = "payments-aks"

default_node_pool {
  name            = "systemnp"
  node_count      = var.node_count_system
  vm_size         = var.node_pool_vm_size
  vnet_subnet_id  = var.aks_subnet_id
  type            = "VirtualMachineScaleSets"  # FIX: only valid value!
  node_labels     = { "nodepool-type" = "system" }
}

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    network_policy    = "calico"
    load_balancer_sku = "standard"
    outbound_type     = "userDefinedRouting"
  }

  private_cluster_enabled     = true
  workload_identity_enabled   = true
  oidc_issuer_enabled        = true

  tags = var.tags
}

resource "azurerm_kubernetes_cluster_node_pool" "usernp" {
  name                  = "usernp"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.aks.id
  vm_size               = var.node_pool_vm_size
  node_count            = var.node_count_user
  mode                  = "User"
  vnet_subnet_id        = var.aks_subnet_id
  node_labels           = { "nodepool-type" = "user" }
  tags                  = var.tags
  orchestrator_version  = "1.28.0"
}
