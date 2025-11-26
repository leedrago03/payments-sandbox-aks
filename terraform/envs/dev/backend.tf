terraform {
  backend "azurerm" {
    resource_group_name  = "rg-payments-aks"
    storage_account_name = ""
    container_name       = ""
    key                  = "dev.terraform.tfstate"
  }
}
