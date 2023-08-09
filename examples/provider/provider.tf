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
