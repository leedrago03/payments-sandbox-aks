output "aks_cluster_name" {
  value = module.aks.aks_cluster_name
}

output "aks_oidc_issuer_url" {
  value = module.aks.oidc_issuer_url
}

output "hub_vnet_id" {
  value = module.networking.hub_vnet_id
}

output "spoke_vnet_id" {
  value = module.networking.spoke_vnet_id
}
