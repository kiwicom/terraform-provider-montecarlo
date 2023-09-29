# Contributing to Terraform provider for Monte Carlo

Welcome to the [terraform-provider-montecarlo](https://github.com/kiwicom/terraform-provider-montecarlo) (by [Kiwi.com](https://kiwi.com)) repository! We're excited that you're interested in contributing to our project.

## Purpose of this Guide

This guide is designed to help contributors understand how to effectively contribute to our project. Whether you're a seasoned developer or new to open source, this guide will provide you with essential information and guidelines on how to make meaningful contributions.

In this guide, you'll find information on:

- [Reporting Issues](#reporting-issues): How to report bugs or request new features.
- [Contributing Code](#contributing-code): Guidelines for making code contributions, including coding standards and PR submission.
- [Documentation](#documentation): Guidelines for documenting your contributions.
- [Licensing](../LICENSE): Information on project licensing.

**Before you get started**, please take a moment to review our [Code of Conduct](#code-of-conduct). We expect all contributors to adhere to our code of conduct, which promotes a positive and inclusive community.

Let's get started on your journey to contributing to [terraform-provider-montecarlo](https://github.com/kiwicom/terraform-provider-montecarlo)!


## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all contributors. To ensure that our community is respectful and supportive, we have adopted a [Code of Conduct](CODE_OF_CONDUCT.md).

Please take a moment to review our Code of Conduct. By participating in this project, you are expected to adhere to its guidelines and help create a positive and inclusive experience for everyone.

If you encounter any behavior that violates our **Code of Conduct**, please report it to the project maintainers by contacting norbert.dopjera@kiwi.com or code@kiwi.com.


## Prerequisites

Contributions to this repository are welcome from anyone who is interested in helping us improve the project. However, to make effective contributions, it's essential to have the following prerequisites:

- **Terraform Knowledge**: Familiarity with _Terraform_ and its usage is crucial.<br><br>
- **Terraform Plugin/Provider Development**: Understanding of _Terraform_ plugin/provider development concepts using **[Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)**.<br><br>
- **GO Programming Language**: Proficiency in the Go programming language, as _Terraform_ providers are primarily written in Go.<br><br>
- **Monte Carlo Knowledge**: Knowledge of the [Monte Carlo](https://www.montecarlodata.com/) data reliability platform and its functionalities.<br><br>
- **GraphQL API**: Understanding of GraphQL APIs is beneficial for contributions since provider communication with **Monte Carlo** is done using its [GraphQL API](https://docs.getmontecarlo.com/docs/using-the-api).<br><br>

While anyone can contribute, having the above prerequisites will help ensure that your contributions are effective and align with the project's goals.

If you're new to Terraform or any of the prerequisites mentioned above, consider exploring our documentation and resources to get started. We appreciate contributions from all skill levels and backgrounds.


## Reporting Issues

We value your contributions, and reporting issues or bugs is an essential part of improving [terraform-provider-montecarlo](https://github.com/kiwicom/terraform-provider-montecarlo). To ensure that issues are clear and actionable, we have provided an [Issue Forms](https://github.com/kiwicom/terraform-provider-montecarlo/issues/new/choose) that you will be asked to choose from when creating a new issue.

**Please follow these guidelines when reporting issues:**

- **Use the [Issue form](https://github.com/kiwicom/terraform-provider-montecarlo/issues/new/choose)**: Our issue form includes fields for essential information, such as a clear title, a detailed description of the issue, steps to reproduce, and the expected and actual behavior.<br><br>
- **Search for Existing Issues**: Before creating a new issue, search our existing issues to see if a similar one has already been reported.<br><br>
- **Be Descriptive**: Provide as much detail as possible in your issue report, including the version of [terraform-provider-montecarlo](https://github.com/kiwicom/terraform-provider-montecarlo), the environment, and any relevant logs or error messages.<br><br>
- **Follow Up**: Be prepared to engage in the conversation around your issue. Our maintainers and contributors may ask for additional details or clarification.<br><br>

We want to make sure that every issue is addressed promptly and effectively. By using our issue forms and following these guidelines, you help us maintain a productive and collaborative environment.

[Create a New Issue](https://github.com/kiwicom/terraform-provider-montecarlo/issues/new/choose)

Thank you for helping us improve [terraform-provider-montecarlo](https://github.com/kiwicom/terraform-provider-montecarlo)!


## Contributing Code

We welcome code contributions to [terraform-provider-montecarlo](https://github.com/kiwicom/terraform-provider-montecarlo)! To ensure that code changes are **consistent and maintainable**, please follow these guidelines:

### Coding Style and Format

Contributors should follow the standard Go programming language coding style and format. You can easily achieve this by using `go fmt` before committing your changes. This repository currently does not yet leverages automatic code style and format checks.

### Submitting Code Changes

- **Pull Requests**: We encourage contributors to submit code changes via Pull Requests (PRs). If you're not familiar with the process, please refer to GitHub's [guide on creating pull requests](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).<br><br>
- **Forking the Repository**: Alternatively, you can fork the repository, make your changes, and create a Pull Request from your fork to [terraform-provider-montecarlo](https://github.com/kiwicom/terraform-provider-montecarlo). This method is suitable if you prefer to work in your own isolated environment.<br><br>

### Branch and Commit Naming

- **Branch Names**: Contributors can **choose** their branch naming strategy, but we encourage the use of [relevant and descriptive](https://tilburgsciencehub.com/building-blocks/collaborate-and-share-your-work/use-github/naming-git-branches/) branch names. If you are planning to persist your branch in long-term fashion, please include any indication of owner (e.g. by prefixing GitHub nickname)<br><br>
- **Commit Messages**: Withing the _Pull requests_ itself and within their own branches, contributors can **choose** their commit naming strategy, but we encourage the use of relevant and descriptive commit names..<br><br>

### Pull Request Workflow

- **Squash and Merge Only**: We use the "squash and merge" option for Pull Requests. This means that each Pull Request will result in a single, well-structured commit on the `master` branch.<br><br>
- **Pull Request Names**: To achieve the required commit message format structure, and to keep our **CI/CD** workflows working, the Pull Request names must follow the `<type>(scope): message` pattern. `type` should be one of the following: `feat`, `fix`, `perf`, `doc`, `test`, `refactor`, `ci`, or any other types allowed by [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).<br><br>
If a commit includes a breaking change, please indicate it in the commit body with a footer `BREAKING CHANGE:`<br><br>

---

These guidelines help us maintain a consistent and manageable codebase. 


## Documentation

Documentation is a crucial part of maintaining and improving [terraform-provider-montecarlo](https://github.com/kiwicom/terraform-provider-montecarlo). Contributors are **required** to document their code changes to ensure that our project remains well-documented and user-friendly.

Providers published to the [Terraform Registry](https://registry.terraform.io/) are required to be uploaded with meaningfull documentation according to registry **[standards](https://developer.hashicorp.com/terraform/registry/providers/docs)**.

### Using `go generate` for Documentation

Contributors can use the `go generate` command, which is automatically set up to create documentation for the **Terraform Registry**. This documentation is sourced from direct Markdown descriptions placed within the schema of Terraform components (such as resources) within the code.

Example of such description can be found directly within the repository. For example in the schema definition of the provider ([provider/provider.go](https://github.com/kiwicom/terraform-provider-montecarlo/blob/master/monte_carlo/provider/provider.go) file).

### Updating the `examples` Folder

To successfuly complete the documentation generation process, contributors are **required** to update the `examples` folder. This folder is used to provide examples in the generated documentation. By keeping examples up-to-date, we provide clear usage examples for users of this provider.

Please ensure that your code changes include appropriate documentation updates, and that examples in the `examples` folder are relevant and accurate.

---

Thank you for helping us maintain comprehensive and user-friendly documentation for this provider!
