# Terraform Provider Monte Carlo

[![GitHub issues](https://img.shields.io/github/issues/kiwicom/terraform-provider-montecarlo)](https://github.com/kiwicom/terraform-provider-montecarlo/issues)
![GitHub go.mod Go version)](https://img.shields.io/github/go-mod/go-version/kiwicom/terraform-provider-montecarlo)
[![last-commit](https://img.shields.io/github/last-commit/kiwicom/terraform-provider-montecarlo)]()
![master cicd status](https://github.com/kiwicom/terraform-provider-montecarlo/actions/workflows/master.yml/badge.svg)
![GitHub release](https://img.shields.io/github/v/release/kiwicom/terraform-provider-montecarlo)
[![Go Report Card](https://goreportcard.com/badge/github.com/kiwicom/terraform-provider-montecarlo)](https://goreportcard.com/report/github.com/kiwicom/terraform-provider-montecarlo)
![coverage](https://raw.githubusercontent.com/kiwicom/terraform-provider-montecarlo/badges/.badges/master/coverage.svg)
![milestone](https://img.shields.io/github/milestones/progress/kiwicom/terraform-provider-montecarlo/2)  

[![Terraform](https://img.shields.io/badge/terraform-%235835CC.svg?style=for-the-badge&logo=terraform&logoColor=white)](https://registry.terraform.io/providers/kiwicom/montecarlo/latest)

The **Terraform Provider Monte Carlo** enables seamless integration with the **[Monte Carlo](https://www.montecarlodata.com/)** data reliability platform, allowing users to **automate** infrastructure provisioning and configuration. With this provider, managing data reliability in your applications becomes effortless, ensuring robust and dependable data pipelines. Simplify your **infrastructure workflows** and enhance your data reliability with this awesome provider.

## Quick Starts

- **Repository rules** - [Code of conduct](./.github/CODE_OF_CONDUCT.md) :memo:
- **For contributors** - [Contributing](./.github/CONTRIBUTING.md) :clipboard:
- **Bug / Feature requests** - [Issue forms](https://github.com/kiwicom/terraform-provider-montecarlo/issues/new/choose) :speech_balloon:

## Installation

To use the **Terraform Provider Monte Carlo**, include it in your Terraform project by adding the following configuration to your `versions.tf` file:

```hcl
terraform {
  required_providers {
    monte_carlo = {
      source  = "kiwicom/montecarlo"
      version = "~> 0.4.0"
    }
  }
}
```

Provider initialization with <ins>**Account Service Key**</ins>, which is used to authenticate API calls of this provider when communicating with **Monte Carlo**.

```hcl
provider "monte_carlo" {
  account_service_key = {
    id    = var.montecarlo_api_key_id     #(secret)
    token = var.montecarlo_api_key_token  #(secret)
  }
}
```

For more information and examples checkout provider documentation either in the [`docs`](docs/index.md) folder or at [Terraform Registry](https://registry.terraform.io/providers/kiwicom/montecarlo/latest/docs).  

This Terraform provider is using **[Protocol Version 6](https://developer.hashicorp.com/terraform/plugin/terraform-plugin-protocol#protocol-version-6)**, making it compatible with Terraform **CLI** of version 1.0 and later. Supported Operating systems and architecture combinations are following:

| _Operating System_  |  `amd64` | `386`  | `arm`  | `arm64`  |
|---:|:---:|:---:|:---:|:---:|
| Windows  | :white_check_mark:  | :white_check_mark:  |   |   |
| Linux  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  |
| Freebsd  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  |   |
| Darwin  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  | :white_check_mark:  |


## Examples

To get started, navigate to the [`examples`](examples/) folder in this repository to find detailed _Terraform_ files for each component provided by the **Terraform Provider Monte Carlo**.  


## License

[![license](https://img.shields.io/github/license/kiwicom/terraform-provider-montecarlo)](https://github.com/kiwicom/terraform-provider-montecarlo/blob/master/LICENSE)
```
MIT License

Copyright (c) 2023 Kiwi.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```