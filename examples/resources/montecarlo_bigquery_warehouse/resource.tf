resource "montecarlo_bigquery_warehouse" "example" {
  name                = "name"
  collector_uuid      = "uuid"
  credentials         = { service_account_key = "{...}" }
  deletion_protection = false
}
