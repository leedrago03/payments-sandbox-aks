variable "resource_group_name" {
  type = string
}

variable "location" {
  type = string
}

variable "redis_name" {
  type = string
}

variable "data_subnet_id" {
  type = string
}

variable "private_dns_zone_id" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
