resource "montecarlo_transactional_warehouse" "example" {
  name           = "name"
  collector_uuid = "uuid"
  db_type        = "POSTGRES" # POSTGRES | MYSQL | SQL-SERVER

  configuration = {
    host     = "host"
    port     = 5432
    database = "database"
    username = "username"  #(secret)
    password = "password"  #(secret)
  }
}
