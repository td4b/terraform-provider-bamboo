terraform {
  required_providers {
    bamboo = {
      source = "td4b/bamboo"
    }
  }
}

provider "bamboo" {
  host    = "http://localhost:8000"
  company = "testcompany"
  apikey  = "APIKEY"
}
