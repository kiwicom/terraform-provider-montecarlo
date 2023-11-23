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
    #"MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++427a1600-2653-40c5-a1e7-5ec98703ee9d++project++gcp-project1-722af1c6",
    #"MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++e7c59fd6-7ca8-41e7-8325-062ea38d3df5++dataset++postgre-dataset-1"
  ]
}
