variable "resource_group_name" {
  type = string
}

variable "location" {
  type = string
}

variable "postgresql_server_name" {
  type = string
}

variable "data_subnet_id" {
  type = string
}

variable "private_dns_zone_id" {
  description = "Private DNS Zone ID for PostgreSQL"
  type        = string
}

variable "admin_username" {
  type    = string
  default = "pgadmin"
}

variable "admin_password" {
  type      = string
  sensitive = true
}

variable "tags" {
  type    = map(string)
  default = {}
}
