terraform {
  required_providers {
    monte_carlo = {
      source  = "kiwicom/montecarlo"
      version = "~> 0.0.1"
    }
  }
}

provider "monte_carlo" {
  account_service_key = {
    id    = "montecarlo"
    token = "montecarlo"
  }
}
