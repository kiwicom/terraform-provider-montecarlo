resource "montecarlo_bigquery_warehouse" "example" {
  name                = "name"
  data_collector_uuid = "uuid"
  service_account_key = "{}"
  deletion_protection = false
}
