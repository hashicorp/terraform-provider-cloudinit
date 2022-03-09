---
layout: "cloudinit"
page_title: "Provider: cloud-init"
description: |-
  The cloud-init provider is used to template strings for other Terraform resources.
---

# Cloud-init Provider

The cloud-init Terraform provider exposes the `cloudinit_config` data source, previously available as the [`template_cloudinit_config` resource in the template provider](https://www.terraform.io/docs/providers/template/d/cloudinit_config.html), which renders a [multipart MIME configuration](https://cloudinit.readthedocs.io/en/latest/topics/format.html#mime-multi-part-archive) for use with [cloud-init](https://cloudinit.readthedocs.io/).

Use the navigation to the left to read about the available data sources.

## Example Usage

```hcl
data "cloudinit_config" "foo" {
  gzip = false
  base64_encode = false

  part {
    content_type = "text/x-shellscript"
    content = "baz"
    filename = "foobar.sh"
  }
}
```
