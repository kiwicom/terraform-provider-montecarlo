terraform {
  required_providers {
    montecarlo = {
      source  = "kiwicom/montecarlo"
      version = "~> 0.2.0"
    }
  }
}

provider "montecarlo" {
  account_service_key = {
    id    = "montecarlo"  #(secret)
    token = "montecarlo"  #(secret)
  }
}
