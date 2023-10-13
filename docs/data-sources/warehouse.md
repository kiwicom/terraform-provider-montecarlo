---
page_title: "montecarlo_warehouse Data Source - terraform-provider-montecarlo"
subcategory: ""
description: |-
  Utility data source which can be used to access MCON's of your data object while using project.dataset.table hierrarchy.
---

# montecarlo_warehouse (Data Source)

Utility _data source_ which can be used to access **MCON's** of your data objects (_assets_) by using `project.dataset.table` hierrarchy. Additionaly, this _data source_ always exports all of the (only **active**) data assets in your warehouse, preventing usage of wrong **MCON's**.

 > _MCON is essentially Monte Carlo universal **identifier** (if you're familiar with AWS you can think of it like the ARN). Its format is following `MCON++{account_uuid}++{resource_uuid}++{object_type}++{object_id}`_  



## Example Usage

```terraform
data "montecarlo_warehouse" "example" {
  uuid = "uuid"
}
```

### Accessing data objects (assets)

```terraform
data "montecarlo_warehouse" "bq" {
  uuid = "427a1600-2653-40c5-a1e7-5ec98703ee9d"
}

resource "montecarlo_domain" "example_assignments_data" {
  name        = "name"
  description = "description"
  assignments = [
    data.montecarlo_warehouse.bq.projects["gcp-project1-722af1c6"].mcon,
    data.montecarlo_warehouse.bq.projects["gcp-project2-744bc2c5"].datasets["postgre-dataset-1"].mcon,
    data.montecarlo_warehouse.bq.projects["gcp-project2-744bc2c5"].datasets["postgre-dataset-2"].tables["table-1"].mcon,
  ]
}
```

Since this _data source_ exposes data objects (assets) **MCON's** in a `map-like` structure (see [below for schema](#schema)), you can use any _Terraform_ utilities for **filtration** and **selection**.

```terraform
data "montecarlo_warehouse" "bq" {
  uuid = "427a1600-2653-40c5-a1e7-5ec98703ee9d"
}

## Assign all datasets from project "gcp-project2-744bc2c5"
## except "postgre-dataset-2" and "postgre-dataset-1"
resource "montecarlo_domain" "example" {
  name        = "name"
  description = "description"
  assignments = values({
    for k, v in montecarlo_warehouse.bq.projects["gcp-project2-744bc2c5"] :
    k => v.mcon
    if !contains(keys(v), "postgre-dataset-2") && !contains(keys(v), "postgre-dataset-1")
  })
}
```



<a id="schema"></a>
## Schema

### Required

- `uuid` (String) Unique identifier of warehouse this resource should collect data objects (assets) from.  

> This resource does not distinguish between different _warehouse_ types (e.g. [bigquery](../resources/bigquery_warehouse.md) vs [transactional](../resources/transactional_warehouse.md)) since all of the _warehouses_ and their data objects (assets) follow the same hierarchy `project.dataset.table`. _Warehouses_ over data sources that are not using any notion of `project` or `dataset` in their hierarchies (for example [transactional warehouse](../resources/transactional_warehouse.md) - `database.schema.table`) alias their hierarchies.

### Read-Only

- `projects` (Attributes Map) All of the projects in the warehouse. (see [below for nested schema](#nestedatt--projects))

<a id="nestedatt--projects"></a>
### Nested Schema for `projects`

Required:

- `mcon` (String) Monte Carlo universal **identifier** of the project.
- `datasets` (Attributes Map) All of the datasets in the project. (see [below for nested schema](#nestedatt--projects--datasets))

<a id="nestedatt--projects--datasets"></a>
### Nested Schema for `projects.datasets`

Required:

- `mcon` (String) Monte Carlo universal **identifier** of the dataset.
- `tables` (Attributes Map) All of the tables in the dataset. (see [below for nested schema](#nestedatt--projects--datasets--tables))

<a id="nestedatt--projects--datasets--tables"></a>
### Nested Schema for `projects.datasets.tables`

Required:

- `mcon` (String) Monte Carlo universal **identifier** of the table.
