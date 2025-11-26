variable "resource_group_name" {
  type = string
}

variable "location" {
  type = string
}

variable "aks_cluster_name" {
  type = string
}

variable "aks_subnet_id" {
  type = string
}

variable "node_pool_vm_size" {
  type = string
}

variable "node_count_system" {
  type = number
}

variable "node_count_user" {
  type = number
}

variable "tags" {
  type    = map(string)
  default = {}
}
