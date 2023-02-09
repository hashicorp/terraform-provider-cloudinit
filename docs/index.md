---
page_title: "cloudinit Provider"
description: |-
  The cloud-init Terraform provider exposes the cloudinit_config data source, previously available as the template_cloudinit_config resource in the template provider https://registry.terraform.io/providers/hashicorp/template/latest/docs/data-sources/cloudinit_config, which renders a multipart MIME configuration https://cloudinit.readthedocs.io/en/latest/explanation/format.html#mime-multi-part-archive for use with cloud-init https://cloudinit.readthedocs.io/en/latest/.
---

# cloudinit Provider

The cloud-init Terraform provider exposes the `cloudinit_config` data source, previously available as the `template_cloudinit_config` resource [in the template provider](https://registry.terraform.io/providers/hashicorp/template/latest/docs/data-sources/cloudinit_config), which renders a [multipart MIME configuration](https://cloudinit.readthedocs.io/en/latest/explanation/format.html#mime-multi-part-archive) for use with [cloud-init](https://cloudinit.readthedocs.io/en/latest/).

This provider requires no configuration. For information on the resources it provides, see the navigation bar.