output "postgresql_server_id" {
  value = azurerm_postgresql_flexible_server.postgres.id
}

output "postgresql_fqdn" {
  value = azurerm_postgresql_flexible_server.postgres.fqdn
}

output "postgresql_database_name" {
  value = azurerm_postgresql_flexible_server_database.payments_db.name
}
