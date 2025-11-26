terraform {
  backend "azurerm" {
    resource_group_name  = "rg-payments-aks"
    storage_account_name = "tfstate9386b294"
    container_name       = "tfstate"
    key                  = "dev.terraform.tfstate"
  }
}
