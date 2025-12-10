output "postgresql_server_id" {
  description = "PostgreSQL Server ID"
  value       = azurerm_postgresql_flexible_server.postgres.id
}

output "postgresql_fqdn" {
  description = "PostgreSQL Server FQDN"
  value       = azurerm_postgresql_flexible_server.postgres.fqdn
}

output "postgresql_database_name" {
  description = "PostgreSQL Database Name"
  value       = azurerm_postgresql_flexible_server_database.payments_db.name
}
