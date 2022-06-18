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

resource "azurerm_resource_group" "three-tier-iac" {
  location = var.location
  name     = "three-tier-rg"
}

resource "azurerm_virtual_network" "three-tier-iac" {
  address_space       = ["10.0.0.0/16"]
  location            = var.location
  name                = "three-tier-network"
  resource_group_name = azurerm_resource_group.three-tier-iac.name
}

