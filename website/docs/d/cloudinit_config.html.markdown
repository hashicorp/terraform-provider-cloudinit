---
layout: "cloudinit"
page_title: "Cloud-init: cloudinit_config"
description: |-
  Renders a multi-part cloud-init config from source files.
---

# cloudinit_config

Renders a [multipart MIME configuration](https://cloudinit.readthedocs.io/en/latest/topics/format.html#mime-multi-part-archive)
for use with [cloud-init](https://cloudinit.readthedocs.io/).

Cloud-init is a commonly-used startup configuration utility for cloud compute
instances. It accepts configuration via provider-specific user data mechanisms,
such as `user_data` for Amazon EC2 instances. Multipart MIME is one of the
data formats it accepts. For more information, see
[User-Data Formats](https://cloudinit.readthedocs.io/en/latest/topics/format.html)
in the cloud-init manual.

This is not a generalized utility for producing multipart MIME messages. Its
featureset is specialized for the features of cloud-init.

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

## Argument Reference

The following arguments are supported:

* `gzip` - (Optional) Specify whether or not to gzip the rendered output. Defaults to `true`.

* `base64_encode` - (Optional) Base64 encoding of the rendered output. Defaults to `true`,
  and cannot be disabled if `gzip` is `true`.

* `boundary` - (Optional) Define the Writer's default boundary separator. Defaults to `MIMEBOUNDARY`.

* `part` - (Required) A nested block type which adds a file to the generated
  cloud-init configuration. Use multiple `part` blocks to specify multiple
  files, which will be included in order of declaration in the final MIME
  document.

Each `part` block expects the following arguments:

* `content` - (Required) Body content for the part.

* `filename` - (Optional) A filename to report in the header for the part.

* `content_type` - (Optional) A MIME-style content type to report in the header for the part.

* `merge_type` - (Optional) A value for the `X-Merge-Type` header of the part,
  to control [cloud-init merging behavior](https://cloudinit.readthedocs.io/en/latest/topics/merging.html).

## Attributes Reference

The following attributes are exported:

* `rendered` - The final rendered multi-part cloud-init config.
