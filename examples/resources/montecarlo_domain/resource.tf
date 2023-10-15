resource "montecarlo_domain" "example_assignments_raw" {
  name        = "name"
  description = "description"
  assignments = [
    "MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++427a1600-2653-40c5-a1e7-5ec98703ee9d++project++gcp-project1-722af1c6",
    "MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++e7c59fd6-7ca8-41e7-8325-062ea38d3df5++dataset++postgre-dataset-1"
  ]
}

resource "montecarlo_domain" "example_tags" {
  name        = "name"
  description = "description"
  tags = [
    {
      name  = "owner"
      value = "bi-team-a"
    },
    {
      name = "montecarlo"
    }
  ]
}

## ITS POSSIBLE TO USE DATA SOURCE WAREHOUSE TO OBTAIN
## ASSIGNMENTS MAPPING TO MCON's AUTOMATICALLY INSTEAD
## OF PROVIDING IT MANUALLY IN RAW FORMAT

data "montecarlo_warehouse" "bq" {
  uuid = "427a1600-2653-40c5-a1e7-5ec98703ee9d"
}

resource "montecarlo_domain" "example_assignments_advanced" {
  name        = "name"
  description = "description"
  assignments = [
    data.montecarlo_warehouse.bq.projects["gcp-project1-722af1c6"].mcon,
    data.montecarlo_warehouse.bq.projects["gcp-project2-744bc2c5"].datasets["postgre-dataset-1"].mcon,
    data.montecarlo_warehouse.bq.projects["gcp-project2-744bc2c5"].datasets["postgre-dataset-2"].tables["table-1"].mcon,
  ]
}
