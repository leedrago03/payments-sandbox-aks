variable "resource_group_name" {
  description = "Name of the Azure resource group"
  type        = string
}

variable "location" {
  description = "Azure region"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "hub_vnet_name" {
  description = "Hub VNet name"
  type        = string
}

variable "spoke_vnet_name" {
  description = "Spoke VNet name"
  type        = string
}

variable "hub_address_space" {
  description = "Hub VNet CIDR block"
  type        = list(string)
}

variable "spoke_address_space" {
  description = "Spoke VNet CIDR block"
  type        = list(string)
}

variable "tags" {
  description = "Tags for all resources"
  type        = map(string)
  default     = {}
}
