provider "azurerm" {
  features {}
}

data "terraform_remote_state" "dev" {
  backend = "azurerm"
  config = {
    resource_group_name  = "rg-payments-aks"
    storage_account_name = "tfstate9386b294"  # Replace with your actual storage account name
    container_name       = "tfstate"
    key                  = "dev.terraform.tfstate"
  }
}

resource "tls_private_key" "jumpbox_ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

module "jumpbox" {
  source                = "../../modules/jumpbox"
  resource_group_name   = "rg-payments-aks"
  location              = "southeastasia"
  vm_name               = "jumpbox-dev"
  vm_size               = "Standard_B1s"
  management_subnet_id  = data.terraform_remote_state.dev.outputs.management_subnet_id
  admin_username        = "azureuser"
  ssh_public_key        = tls_private_key.jumpbox_ssh.public_key_openssh
  aks_cluster_name      = data.terraform_remote_state.dev.outputs.aks_cluster_name
  tags = {
    environment = "dev-tools"
    purpose     = "jumpbox"
  }
}

resource "local_file" "private_key" {
  content         = tls_private_key.jumpbox_ssh.private_key_pem
  filename        = "${path.module}/jumpbox-key.pem"
  file_permission = "0600"
}
