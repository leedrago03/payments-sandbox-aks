variable "resource_group_name" {
  description = "Name of the resource group"
  type        = string
}

variable "location" {
  description = "Azure region"
  type        = string
}

variable "environment" {
  description = "Environment name (dev, prod)"
  type        = string
}

variable "hub_vnet_name" {
  description = "Name of the hub VNet"
  type        = string
}

variable "spoke_vnet_name" {
  description = "Name of the spoke VNet"
  type        = string
}

variable "hub_address_space" {
  description = "Address space for hub VNet"
  type        = list(string)
}

variable "spoke_address_space" {
  description = "Address space for spoke VNet"
  type        = list(string)
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}

