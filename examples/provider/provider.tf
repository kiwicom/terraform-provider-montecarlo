terraform {
  required_providers {
    montecarlo = {
      source  = "kiwicom/montecarlo"
      version = "~> 0.0.1"
    }
  }
}

provider "montecarlo" {
  account_service_key = {
    id    = "montecarlo"
    token = "montecarlo"
  }
}
