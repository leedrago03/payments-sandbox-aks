output "redis_id" {
  value = azurerm_redis_cache.redis.id
}

output "redis_hostname" {
  value = azurerm_redis_cache.redis.hostname
}

output "redis_primary_access_key" {
  value     = azurerm_redis_cache.redis.primary_access_key
  sensitive = true
}
