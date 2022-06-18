resource "azurerm_subnet" "three-tier-logic-iac" {
  address_prefixes     = ["10.0.2.0/24"]
  name                 = "three-tier-logic-subnet"
  resource_group_name  = azurerm_resource_group.three-tier-iac.name
  virtual_network_name = azurerm_virtual_network.three-tier-iac.name
}

resource "azurerm_network_interface" "three-tier-logic-iac" {
  location            = var.location
  name                = "three-tier-logic-nic"
  resource_group_name = azurerm_resource_group.three-tier-iac.name
  ip_configuration {
    name                          = "three-tier-logic-ip"
    private_ip_address_allocation = "Dynamic"
    subnet_id                     = azurerm_subnet.three-tier-logic-iac.id
  }
}

resource "azurerm_linux_virtual_machine" "three-tier-logic-iac" {
  admin_username        = "iacadmin"

  admin_ssh_key {
    public_key = file("~/.ssh/id_rsa.pub")
    username   = "iacadmin"
  }

  location              = var.location
  name                  = "three-tier-logic-machine"

  network_interface_ids = [
    azurerm_network_interface.three-tier-logic-iac.id
  ]

  resource_group_name   = azurerm_resource_group.three-tier-iac.name
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
