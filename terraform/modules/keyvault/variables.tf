variable "resource_group_name" {
  type = string
}

variable "location" {
  type = string
}

variable "keyvault_name" {
  description = "Name of the Key Vault (must be globally unique)"
  type        = string
}

variable "tenant_id" {
  description = "Azure AD Tenant ID"
  type        = string
}

variable "aks_kubelet_identity_object_id" {
  description = "Object ID of AKS kubelet managed identity for RBAC"
  type        = string
}

variable "data_subnet_id" {
  description = "Subnet ID for private endpoint"
  type        = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
