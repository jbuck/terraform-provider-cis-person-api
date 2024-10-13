terraform {
  required_providers {
    cis = {
      source = "registry.terraform.io/mozilla/cis"
    }
  }
}

provider "cis" {}
