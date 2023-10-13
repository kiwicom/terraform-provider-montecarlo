---
page_title: "montecarlo_bigquery_warehouse Resource - terraform-provider-montecarlo"
subcategory: ""
description: |-
  This resource represents the integration of Monte Carlo with BigQuery data warehouse.
---

# montecarlo_bigquery_warehouse (Resource)

Represents the integration of the **Monte Carlo** platform with _BigQuery_ data warehouse. While this resource is not responsible for handling data access and other operations, such as data filtering, it is **responsible for managing the connection** to the _BigQuery_ using the provided service account key.

> **Warning:** Current implementations of _warehouse_ resources are not capable of fetching credentials for given _warehouse_ from the **Monte Carlo** APIs. For this reason its not possible for this resource to detect whether the credentials have been changed externally in the **Monte Carlo**. You can force credentials update by changing remote instance state and re-running your _Terraform_ commands.

To get more information about **Monte Carlo** warehouses, see:
- [API documentation](https://apidocs.getmontecarlo.com/#definition-Warehouse)
- How-to Guides
  - [BigQuery Warehouse](https://docs.getmontecarlo.com/docs/bigquery)



## Example Usage

```terraform
resource "montecarlo_bigquery_warehouse" "example" {
  name                = "name"
  collector_uuid      = "uuid"
  service_account_key = "{...}"
}
```



## Schema

### Required

- `name` (String) Name of the warehouse instance, as it should be presented in the **Monte Carlo**.  

- `collector_uuid` (String) Unique identifier of data collector this warehouse instance will be attached to. You can find all of your data collectors in the **Monte Carlo** _Settings_ -> _Integrations_ -> _Collectors_ page.  

  - If changed in the _Terraform_ configuration, resource instance will be **deleted** (leading to a new resource creation on the next `terraform plan/apply`).  

  - If changed in the remote instance state, resource instance will be **removed** from the _Terraform_ state, but not deleted (leading to a new resource creation on the next `terraform plan/apply`).  

- `service_account_key` (String, Sensitive) Service account key used by the warehouse connection for authentication and authorization against _BigQuery_. The very same service account is used to grant required permissions to _Monte Carlo BigQuery warehouse_ for the data access. For more information follow **Monte Carlo** documentation: https://docs.getmontecarlo.com/docs/bigquery.  

  - The key must be provided as raw **JSON** string (_as obtained when downloaded from GCP UI_), otherwise this reasource will fail during _Terraform_ commands.  

### Optional

- `deletion_protection` (Boolean, _default:_ `true`) Unless this field is set to false, a terraform destroy or terraform apply that would delete the instance **will fail**, leaving the instance unchanged. This setting will prevent the deletion even if the resource instance is already deleted.

### Read-Only

- `uuid` (String) Unique identifier of warehouse managed by this resource.  

- `connection_uuid` (String) Unique identifier of connection managed by this resource, responsible for communication with _BigQuery_.  

  - if _connection type_ of the connection managed by this reasource changes externally, this reasource **will fail to read** external state (_blocking any further resource functionality_). In such scenario a manual intervention is required.



## Import

This resource can be imported using the import ID with following format:

* `{{<warehouse_uuid>,<connection_uuid>,<data_collector_uuid>}}`

In **Terraform v1.5.0** and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import a _BigQuery Warehouse_ using one of the formats above. For example:

```terraform
import {
  id = "{{importID}}"
  to = montecarlo_bigquery_warehouse.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), _BigQuery Warehouse_ can be imported using one of the formats above. For example:

```
$ terraform import montecarlo_bigquery_warehouse.default {{importID}}
```
