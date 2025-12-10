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

variable "oidc_issuer_url" {
  description = "AKS OIDC issuer URL"
  type        = string
}

variable "eventhubs_namespace_id" {
  description = "Event Hubs Namespace ID for RBAC"
  type        = string
}

variable "keyvault_id" {
  description = "Key Vault ID for RBAC"
  type        = string
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default     = {}
}
