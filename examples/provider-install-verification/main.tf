terraform {
  required_providers {
    bamboo = {
      source = "hashicorp.com/edu/bamboo"
    }

  }
}

provider "bamboo" {
  host    = "http://localhost:8000"
  company = "testcompany"
  apikey  = "APIKEY"
}

data "bamboo_users" "bambootest" {}

output "result" {
  value = data.bamboo_users.bambootest.users
}
