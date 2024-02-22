

<a name="v0.4.0"></a>
## v0.4.0

> 2024-02-22

- Full diff - **[v0.3.0...v0.4.0](https://github.com/kiwicom/terraform-provider-montecarlo/compare/v0.3.0...v0.4.0)**  

### :books: Documentation (unchanged functionality)

* **warehouse:** drift credentials examples and docs ([#89](https://github.com/kiwicom/terraform-provider-montecarlo/issues/89))

### :bug: Bug Fixes

* **warehouse:** transactional credentials drift changes ([#85](https://github.com/kiwicom/terraform-provider-montecarlo/issues/85))
* **warehouse:** BQ credentials drift changes ([#84](https://github.com/kiwicom/terraform-provider-montecarlo/issues/84))

### :mag: Tests (unchanged functionality)

* **all:** acc. test working against real infrastructure ([#68](https://github.com/kiwicom/terraform-provider-montecarlo/issues/68))

### :sparkles: Features

* **monitors:** comparison_monitor basic create - stopped development ([#82](https://github.com/kiwicom/terraform-provider-montecarlo/issues/82))
* **warehouse:** connection handling refactoring ([#86](https://github.com/kiwicom/terraform-provider-montecarlo/issues/86))


<a name="v0.3.0"></a>
## v0.3.0

> 2023-11-01

- Full diff - **[v0.2.1...v0.3.0](https://github.com/kiwicom/terraform-provider-montecarlo/compare/v0.2.1...v0.3.0)**  

### :sparkles: Features

* **iam:** IAM member resource implemented ([#53](https://github.com/kiwicom/terraform-provider-montecarlo/issues/53))
* **iam:** authorization groups resource implementation ([#50](https://github.com/kiwicom/terraform-provider-montecarlo/issues/50))

### :bug: Bug Fixes

* **iam:** iam_member using groups API for assignment ([#61](https://github.com/kiwicom/terraform-provider-montecarlo/pull/61))

### :books: Documentation (unchanged functionality)

* **gen:** removed attributes for docs generation ([#51](https://github.com/kiwicom/terraform-provider-montecarlo/issues/51))
* **resources:** iam_member documentation and examples ([#63](https://github.com/kiwicom/terraform-provider-montecarlo/issues/63))
* **resources:** iam_group documentation and examples ([#60](https://github.com/kiwicom/terraform-provider-montecarlo/issues/60))

### :mag: Tests (unchanged functionality)

* **iam:** member assignment acceptance tests ([#55](https://github.com/kiwicom/terraform-provider-montecarlo/issues/55))
* **iam:** IAM group resource acceptance tests ([#54](https://github.com/kiwicom/terraform-provider-montecarlo/issues/54))


<a name="v0.2.1"></a>
## v0.2.1

> 2023-10-17

- Full diff - **[v0.2.0...v0.2.1](https://github.com/kiwicom/terraform-provider-montecarlo/compare/v0.2.0...v0.2.1)**  

### :bug: Bug Fixes

* **resources:** bigquery warehouse missing state upgrade v0 ([#48](https://github.com/kiwicom/terraform-provider-montecarlo/issues/48))


<a name="v0.2.0"></a>
## v0.2.0

> 2023-10-16

- Full diff - **[v0.2.0-pre...v0.2.0](https://github.com/kiwicom/terraform-provider-montecarlo/compare/v0.2.0-pre...v0.2.0)**  

### :books: Documentation (unchanged functionality)

* **all:** add missing docs before 0.2.0 release ([#45](https://github.com/kiwicom/terraform-provider-montecarlo/issues/45))

### :bug: Bug Fixes

* **resources:** resolved common issues before 0.2.0 ([#44](https://github.com/kiwicom/terraform-provider-montecarlo/issues/44))


<a name="v0.2.0-pre"></a>
## v0.2.0-pre

> 2023-10-11

- Full diff - **[v0.1.3...v0.2.0-pre](https://github.com/kiwicom/terraform-provider-montecarlo/compare/v0.1.3...v0.2.0-pre)**  

### :sparkles: Features

* **data_sources:** warehouse exposing MCON's implementation ([#41](https://github.com/kiwicom/terraform-provider-montecarlo/issues/41))
* **resource|postgres_warehouse:** kick-off beta version implementation ([#40](https://github.com/kiwicom/terraform-provider-montecarlo/issues/40))


<a name="v0.1.3"></a>
## v0.1.3

> 2023-10-09

- Full diff - **[v0.1.2...v0.1.3](https://github.com/kiwicom/terraform-provider-montecarlo/compare/v0.1.2...v0.1.3)**  

### :books: Documentation (unchanged functionality)

* **resource|domain:** added docs and terraform examples ([#38](https://github.com/kiwicom/terraform-provider-montecarlo/issues/38))

### :sparkles: Features

* **resource|domain:** import state via domain uuid ([#37](https://github.com/kiwicom/terraform-provider-montecarlo/issues/37))
* **resource|domain:** kick-off beta version implementation ([#35](https://github.com/kiwicom/terraform-provider-montecarlo/issues/35))


<a name="v0.1.2"></a>
## v0.1.2

> 2023-10-04

- Full diff - **[v0.1.1...v0.1.2](https://github.com/kiwicom/terraform-provider-montecarlo/compare/v0.1.1...v0.1.2)**  

### :books: Documentation (unchanged functionality)

* **examples:** fixed provider instalation syntax and usage ([#33](https://github.com/kiwicom/terraform-provider-montecarlo/issues/33))

### :bug: Bug Fixes

* **warehouse|bigquery:** read operation - response missing dataCollector ([#32](https://github.com/kiwicom/terraform-provider-montecarlo/issues/32))


<a name="v0.1.1"></a>
## v0.1.1

> 2023-10-04

- Full diff - **[v0.1.0...v0.1.1](https://github.com/kiwicom/terraform-provider-montecarlo/compare/v0.1.0...v0.1.1)**  

### :bug: Bug Fixes

* **warehouse|bigquery:** read operation - inconsistent data collector ([#29](https://github.com/kiwicom/terraform-provider-montecarlo/issues/29))


<a name="v0.1.0"></a>
## [v0.1.0](https://github.com/kiwicom/terraform-provider-montecarlo/compare/v0.0.1...v0.1.0)

> 2023-09-30

### Bug Fixes

* **release|registry:** provider naming and docs for registry ([#8](https://github.com/kiwicom/terraform-provider-montecarlo/issues/8))
* **repo|reports:** broken links and unnecessary newlines ([#13](https://github.com/kiwicom/terraform-provider-montecarlo/issues/13))

### Features

* **release|registry:** manifest for Terraform registry ([#9](https://github.com/kiwicom/terraform-provider-montecarlo/issues/9))
* **repo:** github changelog and release notes ([#10](https://github.com/kiwicom/terraform-provider-montecarlo/issues/10))
* **repo|docs:** finished basic community standard documents ([#11](https://github.com/kiwicom/terraform-provider-montecarlo/issues/11))
* **repo|reports:** migrate to beta issue forms ([#12](https://github.com/kiwicom/terraform-provider-montecarlo/issues/12))
* **resources:** import state for 'bigquery_warehouse' ([#14](https://github.com/kiwicom/terraform-provider-montecarlo/issues/14))


<a name="v0.0.1"></a>
## v0.0.1

> 2023-09-27

### Bug Fixes

* **ci-cd|release:** initial goreleaser configuration syntax error ([#7](https://github.com/kiwicom/terraform-provider-montecarlo/issues/7))

### Features

* **ci-cd|release:** initial configuration for golang binaries ([#6](https://github.com/kiwicom/terraform-provider-montecarlo/issues/6))
* **provider:** added Code of Conduct ([#5](https://github.com/kiwicom/terraform-provider-montecarlo/issues/5))
* **provider:** golang initialization and BigQuery warehouse implementation  ([#2](https://github.com/kiwicom/terraform-provider-montecarlo/issues/2))
* **provider|code:** added MIT license ([#4](https://github.com/kiwicom/terraform-provider-montecarlo/issues/4))
* **resources|warehouses:** data_collector_uuid attribute exposed ([#3](https://github.com/kiwicom/terraform-provider-montecarlo/issues/3))

