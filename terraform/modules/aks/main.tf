resource "azurerm_kubernetes_cluster" "aks" {
  name                = var.aks_cluster_name
  location            = var.location
  resource_group_name = var.resource_group_name
  dns_prefix          = "payments-aks"
  kubernetes_version  = "1.32.9"

  default_node_pool {
    name                = "systemnp"
    node_count          = var.node_count_system
    vm_size             = var.system_node_pool_vm_size
    vnet_subnet_id      = var.aks_subnet_id
    type                = "VirtualMachineScaleSets"
    node_labels = {
      "nodepool-type" = "system"
      "workload-type" = "infrastructure"
    }
    # node_taints removed - not supported in default_node_pool
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin    = "azure"
    network_policy    = "calico"
    load_balancer_sku = "standard"
    outbound_type     = "loadBalancer"
  }

  private_cluster_enabled     = true
  workload_identity_enabled   = true
  oidc_issuer_enabled        = true

  tags = var.tags
}

resource "azurerm_kubernetes_cluster_node_pool" "usernp" {
  name                  = "usernp"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.aks.id
  vm_size               = var.user_node_pool_vm_size
  node_count            = var.node_count_user
  mode                  = "User"
  vnet_subnet_id        = var.aks_subnet_id
  node_labels = {
    "nodepool-type" = "user"
    "workload-type" = "application"
  }
  tags = var.tags
}
