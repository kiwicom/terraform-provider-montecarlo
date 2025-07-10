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

variable "bq_service_account" {
  type = string
}

resource "montecarlo_bigquery_warehouse" "test" {
  name                = "test-warehouse-updated"
  collector_uuid      = "9d1aee0a-6a90-47f0-8221-a884be707fc4"
  credentials         = { service_account_key = var.bq_service_account }
  deletion_protection = false
}
