output "jumpbox_ip" {
  value = module.jumpbox.jumpbox_public_ip
}

output "ssh_command" {
  value = module.jumpbox.ssh_command
}

output "connect_instructions" {
  value = <<-EOT
    1. SSH to jumpbox: ssh -i jumpbox-key.pem azureuser@${module.jumpbox.jumpbox_public_ip}
    2. Once logged in, run: ./connect-aks.sh
    3. Test with: kubectl get nodes
  EOT
}
