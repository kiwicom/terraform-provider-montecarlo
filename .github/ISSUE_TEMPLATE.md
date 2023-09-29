---
name: Bug Report / Feature Request
about: Create a report to help us improve the Terraform Provider
---

### Issue Type

- [ ] Bug Report
- [ ] Feature Request

If this is an **Feature Request**, feel free to skip following sections and jump ahead to [Feature Request Details](#feature-request-details). If you feel like any of the following sections is not related to this report, do not include it.

### Actual Behavior

Explain what actually happened. Include any error messages or logs if applicable. You are free to include **Terraform plan** properly formatted outputs here.

### Steps to Reproduce

1. Provide a step-by-step guide to reproduce the issue.
2. Include Terraform configurations and any necessary resources.
3. If possible, provide a minimal, complete, and verifiable example (MCVE).

### Terraform Information

 - Please specify the version of Terraform you are using (e.g., 0.15.0):
 - Please specify the version of the Terraform Provider you are using (e.g., v0.1.0):

### Environment Information

- Operating System: [e.g., Windows 10, macOS, Linux]
- Relevant Environment Variables or Settings: [if applicable]

### Terraform Configuration

- list all of the affected resources here 
- montecarlo_xxx

_Copy-paste your Terraform configurations here. Use Markdown code block with propper formatting. If reproducing the bug involves modifying the config file (e.g., apply a config, change a value, apply the config again, see the bug) then please include both the version of the config before the change, and the version of the config after the change._

```hcl
# Replace sensitive information with placeholders.
resource "montecarlo_xxx" "name" {
  # Configuration options
}
```

### Related Issues

Are there any existing issues related to this one? Please link them here.

### Checklist

- [ ] I have reviewed the documentation and other relevant resources.
- [ ] I have searched for similar issues in the repository.
- [ ] I have provided all the requested information.
- [ ] I have tested the issue with the latest version of the Terraform Provider (if possible).

### Feature Request Details

If this is a feature request, please describe the new functionality you would like to see.

### Acceptance Criteria (for Feature Requests)

If this is a feature request, specify the acceptance criteria for the proposed feature.

### Labels

Please label this issue appropriately:

- [ ] Bug
- [ ] Enhancement (for feature requests)
- [ ] Needs Triage
- [ ] Help Wanted

### Additional Information

Any additional information or context you would like to provide. For example:

- Relevant logs or error output.
- Relevant configuration files or code snippets.
- Any potential workarounds you've tried.
- Links or any sources you found yourself while facing the issue

---

Thank you for contributing to our Terraform Provider repository!
