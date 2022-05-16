terraform {
  required_version = ">=0.12"

  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "~>2.0"
    }
  }
}

provider "azurerm" {
  features {}
}

variable "location" {
  description = "Azure region to deploy to."
  default = "Central US"
}

resource "azurerm_resource_group" "vm-iac" {
  name = "vm-iac-example"
  location = var.location
}

resource "azurerm_virtual_network" "vm-iac" {
  address_space       = ["10.0.0.0/16"]
  location            = var.location
  name                = "vm-iac-network"
  resource_group_name = azurerm_resource_group.vm-iac.name
}

resource "azurerm_subnet" "vm-iac" {
  address_prefixes     = ["10.0.2.0/24"]
  name                 = "vm-iac-subnet"
  resource_group_name  = azurerm_resource_group.vm-iac.name
  virtual_network_name = azurerm_virtual_network.vm-iac.name
}

resource "azurerm_public_ip" "vm-iac" {
  allocation_method   = "Dynamic"
  location            = var.location
  name                = "vm-iac-public-ip"
  resource_group_name = azurerm_resource_group.vm-iac.name
}

resource "azurerm_network_interface" "vm-iac" {
  location            = var.location
  name                = "vm-iac-nic"
  resource_group_name = azurerm_resource_group.vm-iac.name
  ip_configuration {
    name                          = "vm-ip"
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.vm-iac.id
    subnet_id                     = azurerm_subnet.vm-iac.id
  }
}

resource "azurerm_network_security_group" "vm-iac" {
  location            = var.location
  name                = "vm-iac-nsg"
  resource_group_name = azurerm_resource_group.vm-iac.name

  security_rule {
    access    = "Allow"
    destination_address_prefix = "*"
    destination_port_range = "22"
    direction = "Inbound"
    name      = "vm-iac-ssh-access-rule"
    priority  = 100
    protocol  = "TCP"
    source_address_prefix = "*"
    source_port_range = "*"
  }
}

resource "azurerm_subnet_network_security_group_association" "vm-iac" {
  subnet_id = azurerm_subnet.vm-iac.id
  network_security_group_id = azurerm_network_security_group.vm-iac.id
}

resource "azurerm_linux_virtual_machine" "vm-iac" {
  admin_username        = "iacadmin"

  admin_ssh_key {
    username   = "iacadmin"
    public_key = file("~/.ssh/id_rsa.pub")
  }

  location              = var.location
  name                  = "vm-iac-machine"

  network_interface_ids = [
    azurerm_network_interface.vm-iac.id
  ]

  resource_group_name   = azurerm_resource_group.vm-iac.name
  size                  = "Standard_F2"

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    offer     = "UbuntuServer"
    publisher = "Canonical"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}