# This resource is currently only in Beta version ! Assignments such as projects,
# datasets and tables can only be provided in MCON's format (Monte Carlo universal identifiers).
# Stable release for this resource will contain automatic translation of data assets to the MCON's
#
# Users using `tags` attribute instead of assignments can consider this functionality as stable
# MCON format is following: MCON++{account_uuid}++{resource_uuid}++{object_type}++{object_id}

resource "montecarlo_domain" "example_assignments" {
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
