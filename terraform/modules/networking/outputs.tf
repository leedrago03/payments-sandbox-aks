output "hub_vnet_id" {
  value = azurerm_virtual_network.hub.id
}
output "spoke_vnet_id" {
  value = azurerm_virtual_network.spoke.id
}
output "aks_subnet_id" {
  value = azurerm_subnet.aks.id
}
output "ingress_subnet_id" {
  value = azurerm_subnet.ingress.id
}
output "data_subnet_id" {
  value = azurerm_subnet.data.id
}
output "postgres_subnet_id" {
  value = azurerm_subnet.postgres.id
}

output "management_subnet_id" {
  value = azurerm_subnet.management.id
}
