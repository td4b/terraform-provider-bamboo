terraform {
  required_providers {
    bamboo = {
      source = "github.com/td4b/terraform-provider-bamboo"
    }

  }
}

provider "bamboo" {
  host    = "http://localhost:8000"
  company = "testcompany"
  apikey  = "APIKEY"
}

data "bamboo_users" "bambootest" {}

module "ad_test" {
  for_each = data.bamboo_users.bambootest
  source   = "./test-module"
  values   = each.value
}

# output "result" {
#   value = data.bamboo_users.bambootest
# }
