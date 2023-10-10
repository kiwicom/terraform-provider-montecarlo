---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "montecarlo_domain Resource - terraform-provider-montecarlo"
subcategory: ""
description: |-
  (Beta version !!) A named resource which lets you define a collection of tables or views by selecting a combination of tables, schemas or databases. Domains can be used to create notifications and authorization groups as a way to adjust the scope without having to redefine a list of tables every time.
---

# montecarlo_domain (Resource)

**(Beta version !!)** A named resource which lets you define a collection of tables or views by selecting a combination of tables, schemas or databases. Domains can be used to create notifications and authorization groups as a way to adjust the scope without having to redefine a list of tables every time.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Domain as it will be presented in Monte Carlo.

### Optional

- `assignments` (List of String) Objects assigned to domain (in MCONs format: MCON++{account_uuid}++{resource_uuid}++{object_type}++{object_id}).
- `description` (String) Description of the domain as it will be presented in Monte Carlo.
- `tags` (Attributes List) Filter by tag key/value pairs for tables. (see [below for nested schema](#nestedatt--tags))

### Read-Only

- `uuid` (String) Unique identifier of domain managed by this resource.

<a id="nestedatt--tags"></a>
### Nested Schema for `tags`

Required:

- `name` (String) Tag name

Optional:

- `value` (String) Tag value