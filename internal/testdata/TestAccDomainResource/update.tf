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
  description = "Domain test description"
  tags        = [
	{
	  name = "dataset_tables_2"
	}
  ]
}
