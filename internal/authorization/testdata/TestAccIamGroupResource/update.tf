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

resource "montecarlo_iam_group" "test" {
	name = "group-1"
	role = "mcd/viewer"
    sso_group = "ssoGroup1"
    domains = [
        "ba0c4080-089d-4377-8878-466c31d19807",
        "dd4cda19-1c5c-4339-9628-76376c9e281e"
    ]
}
