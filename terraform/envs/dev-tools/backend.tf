terraform {
  backend "azurerm" {
    resource_group_name  = "rg-payments-aks"
    storage_account_name = "tfstate9386b294"  # Replace with your actual storage account name
    container_name       = "tfstate"
    key                  = "dev-tools.terraform.tfstate"
  }
}
