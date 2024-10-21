terraform {
  required_providers {
    cis = {
      source = "registry.terraform.io/mozilla/cis"
    }
  }
}

provider "cis" {}

data "cis_people" "example" {
  email = "jbuckley@mozilla.com"
}

output "example" {
  value = data.cis_people.example
}
