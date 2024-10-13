terraform {
  required_providers {
    cis = {
      source = "registry.terraform.io/mozilla/cis"
    }
  }
}

provider "cis" {}

data "cis_example" "example" {
  configurable_attribute = "some-value"
}
