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
  name        = "domain1"
  description = "Domain test description"
  tags        = [
	{
	  name = "owner"
	  value = "bi-internal"
	},
	{
	  name = "dataset_tables_1"
	}
  ]
}
