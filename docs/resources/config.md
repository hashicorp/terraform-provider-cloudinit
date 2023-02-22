---
page_title: "cloudinit_config Resource - terraform-provider-cloudinit"
description: |-
  Renders a multi-part MIME configuration https://cloudinit.readthedocs.io/en/latest/explanation/format.html#mime-multi-part-archive for use with cloud-init https://cloudinit.readthedocs.io/en/latest/.
  Cloud-init is a commonly-used startup configuration utility for cloud compute instances. It accepts configuration via provider-specific user data mechanisms, such as user_data for Amazon EC2 instances. Multi-part MIME is one of the data formats it accepts. For more information, see User-Data Formats https://cloudinit.readthedocs.io/en/latest/explanation/format.html in the cloud-init manual.
  This is not a generalized utility for producing multi-part MIME messages. It's feature set is specialized for cloud-init multi-part MIME messages.
---

# cloudinit_config (Resource)

~> **This resource is deprecated** Please use the [cloudinit_config](https://registry.terraform.io/providers/hashicorp/cloudinit/latest/docs/data-sources/config)
  data source instead.

Renders a [multi-part MIME configuration](https://cloudinit.readthedocs.io/en/latest/explanation/format.html#mime-multi-part-archive) for use with [cloud-init](https://cloudinit.readthedocs.io/en/latest/).

Cloud-init is a commonly-used startup configuration utility for cloud compute instances. It accepts configuration via provider-specific user data mechanisms, such as `user_data` for Amazon EC2 instances. Multi-part MIME is one of the data formats it accepts. For more information, see [User-Data Formats](https://cloudinit.readthedocs.io/en/latest/explanation/format.html) in the cloud-init manual.

This is not a generalized utility for producing multi-part MIME messages. It's feature set is specialized for cloud-init multi-part MIME messages.

## Example Usage

### Config
```terraform
resource "cloudinit_config" "foobar" {
  gzip          = false
  base64_encode = false

  part {
    filename     = "hello-script.sh"
    content_type = "text/x-shellscript"

    content = file("${path.module}/hello-script.sh")
  }

  part {
    filename     = "cloud-config.yaml"
    content_type = "text/cloud-config"

    content = file("${path.module}/cloud-config.yaml")
  }
}
```

### hello-script.sh
```shell
#!/bin/sh
echo "Hello World! I'm starting up now at $(date -R)!"
```

### cloud-config.yaml
```yaml
#cloud-config
# See documentation for more configuration examples
# https://cloudinit.readthedocs.io/en/latest/reference/examples.html 

# Install arbitrary packages
# https://cloudinit.readthedocs.io/en/latest/reference/examples.html#install-arbitrary-packages
packages:
  - python
# Run commands on first boot
# https://cloudinit.readthedocs.io/en/latest/reference/examples.html#run-commands-on-first-boot
runcmd:
 - [ ls, -l, / ]
 - [ sh, -xc, "echo $(date) ': hello world!'" ]
 - [ sh, -c, echo "=========hello world=========" ]
 - ls -l /root
```

<!-- This schema was originally generated with tfplugindocs, then modified manually to ensure `part` block list is noted as Required -->

## Schema

### Required

- `part` (Block List) A nested block type which adds a file to the generated cloud-init configuration. Use multiple `part` blocks to specify multiple files, which will be included in order of declaration in the final MIME document. (see [below for nested schema](#nestedblock--part))

### Optional

- `base64_encode` (Boolean) Specify whether or not to base64 encode the `rendered` output. Defaults to `true`, and cannot be disabled if gzip is `true`.
- `boundary` (String) Specify the Writer's default boundary separator. Defaults to `MIMEBOUNDARY`.
- `gzip` (Boolean) Specify whether or not to gzip the `rendered` output. Defaults to `true`.

### Read-Only

- `id` (String) [CRC-32](https://pkg.go.dev/hash/crc32) checksum of `rendered` cloud-init config.
- `rendered` (String) The final rendered multi-part cloud-init config.

<a id="nestedblock--part"></a>
### Nested Schema for `part`

Required:

- `content` (String) Body content for the part.

Optional:

- `content_type` (String) A MIME-style content type to report in the header for the part. Defaults to `text/plain`
- `filename` (String) A filename to report in the header for the part.
- `merge_type` (String) A value for the `X-Merge-Type` header of the part, to control [cloud-init merging behavior](https://cloudinit.readthedocs.io/en/latest/reference/merging.html).
