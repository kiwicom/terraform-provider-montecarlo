terraform {
  required_providers {
    montecarlo = {
      source  = "kiwicom/montecarlo"
      version = "~> 0.5.0"
    }
  }
}

provider "montecarlo" {
  account_service_key = {
    id    = var.montecarlo_api_key_id     # (secret)
    token = var.montecarlo_api_key_token  # (secret)
  }
}
