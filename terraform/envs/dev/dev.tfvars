location = "southeastasia"
environment = "dev"
resource_group_name = "rg-payments-aks"
hub_vnet_name = "vnet-hub-dev"
spoke_vnet_name = "vnet-spoke-dev"
hub_address_space = ["10.0.0.0/16"]
spoke_address_space = ["10.1.0.0/16"]
tags = {
  project = "payments-sandbox"
  environment = "dev"
}
