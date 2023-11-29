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

variable "pg_host" {
  type = string
}

variable "pg_port" {
  type = number
}

variable "pg_database" {
  type = string
}

variable "pg_user" {
  type = string
}

variable "pg_password" {
  type = string
}

resource "montecarlo_transactional_warehouse" "test" {
  name                = "name1"
  collector_uuid      = "a08d23fc-00a0-4c36-b568-82e9d0e67ad8"
  db_type             = "POSTGRES" # POSTGRES | MYSQL | SQL-SERVER
  deletion_protection = false

  configuration = {
    host     = var.pg_host
    port     = var.pg_port
    database = var.pg_database
    username = var.pg_user      #(secret)
    password = var.pg_password  #(secret)
  }
}
