terraform {
  required_providers {
    csi = {
      source = "registry.terraform.io/mozilla/cis"
    }
  }
}

provider "csi" {}
