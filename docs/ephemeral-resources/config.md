---
page_title: "cloudinit_config Ephemeral Resource - terraform-provider-cloudinit"
description: |-
  Renders a multi-part MIME configuration https://cloudinit.readthedocs.io/en/latest/explanation/format.html#mime-multi-part-archive for use with cloud-init https://cloudinit.readthedocs.io/en/latest/.
  Cloud-init is a commonly-used startup configuration utility for cloud compute instances. It accepts configuration via provider-specific user data mechanisms, such as user_data for Amazon EC2 instances. Multi-part MIME is one of the data formats it accepts. For more information, see User-Data Formats https://cloudinit.readthedocs.io/en/latest/explanation/format.html in the cloud-init manual.
  This is not a generalized utility for producing multi-part MIME messages. Its feature set is specialized for cloud-init multi-part MIME messages.
---

# cloudinit_config (Ephemeral Resource)

Renders a [multi-part MIME configuration](https://cloudinit.readthedocs.io/en/latest/explanation/format.html#mime-multi-part-archive) for use with [cloud-init](https://cloudinit.readthedocs.io/en/latest/).

Cloud-init is a commonly-used startup configuration utility for cloud compute instances. It accepts configuration via provider-specific user data mechanisms, such as `user_data` for Amazon EC2 instances. Multi-part MIME is one of the data formats it accepts. For more information, see [User-Data Formats](https://cloudinit.readthedocs.io/en/latest/explanation/format.html) in the cloud-init manual.

This is not a generalized utility for producing multi-part MIME messages. Its feature set is specialized for cloud-init multi-part MIME messages.

**This ephemeral resource supports ephemeral values** (such as secrets from Vault KV v2) in the `content` attribute, allowing secrets to be injected directly into cloud-init templates **without storing them in Terraform state**. This enables fully declarative, secret-safe cloud-init generation entirely within Terraform.

~> **Note:** Ephemeral resources are re-evaluated on every plan/apply. When using the `rendered` output with resources like `aws_instance.user_data`, set `user_data_replace_on_change = false` to prevent instance replacement on every run.

## Example Usage

### Basic Usage

#### Config
```terraform
ephemeral "cloudinit_config" "foobar" {
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

#### hello-script.sh
```shell
#!/bin/sh
echo "Hello World! I'm starting up now at $(date -R)!"
```

#### cloud-config.yaml
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

### Usage with AWS Instance

When using with `aws_instance`, set `user_data_replace_on_change = false` to prevent instance replacement:

```terraform
ephemeral "cloudinit_config" "example" {
  part {
    content_type = "text/cloud-config"
    content      = ephemeral.vault_kv_secret_v2.example.data["cloud_config"]
  }
}

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
  user_data     = ephemeral.cloudinit_config.example.rendered
  
  # Prevent instance replacement when ephemeral values change
  user_data_replace_on_change = false
}
```

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

- `content` (String) Body content for the part. **Supports ephemeral values from providers like Vault.**

Optional:

- `content_type` (String) A MIME-style content type to report in the header for the part. Defaults to `text/plain`
- `filename` (String) A filename to report in the header for the part.
- `merge_type` (String) A value for the `X-Merge-Type` header of the part, to control [cloud-init merging behavior](https://cloudinit.readthedocs.io/en/latest/reference/merging.html).