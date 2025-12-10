variable "resource_group_name" {
  description = "Name of the resource group"
  type        = string
}

variable "location" {
  description = "Azure region"
  type        = string
}

variable "environment" {
  description = "Environment name (dev/prod)"
  type        = string
}

variable "data_subnet_id" {
  description = "Subnet ID for private endpoint"
  type        = string
}

variable "vnet_id" {
  description = "Virtual Network ID for DNS zone linking"
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
