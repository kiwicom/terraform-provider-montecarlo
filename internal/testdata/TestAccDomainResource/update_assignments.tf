variable "montecarlo_api_key_id" {
  type = string
}

variable "montecarlo_api_key_token" {
  type = string
}

provider "montecarlo" {
  account_service_key = {
    id    = var.montecarlo_api_key_id     # (secret)
    token = var.montecarlo_api_key_token  # (secret)
  }
}

resource "montecarlo_domain" "test" {
  name        = "domain2"
  description = "Domain test description 2"
  assignments = [
    "MCON++3e9abc75-5dc1-447e-b4cb-9d5a6fc5db5c++da6c0716-2724-4bfc-b5cc-7e0364faf979++table++data-playground-8bb9fc23:terraform_provider_montecarlo.person",
    "MCON++3e9abc75-5dc1-447e-b4cb-9d5a6fc5db5c++da6c0716-2724-4bfc-b5cc-7e0364faf979++table++data-playground-8bb9fc23:terraform_provider_montecarlo.device"
  ]
}
