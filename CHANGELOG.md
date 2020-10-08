## 2.0.0 (Unreleased)

Binary releases of this provider will now include the linux-arm64 platform.

BREAKING CHANGES:

* Upgrade to version 2 of the Terraform Plugin SDK, which drops support for Terraform 0.11. This provider will continue to work as expected for users of Terraform 0.11, which will not download the new version. [GH-3]

## 1.0.0 (April 14, 2020)

Initial release. This provider exposes one resource, cloudinit_config, which is identical to the template_cloudinit_config resource in terraform-provider-template.
