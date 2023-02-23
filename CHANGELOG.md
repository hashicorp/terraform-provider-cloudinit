## 2.3.2 (February 23, 2023)

BUG FIXES:

* cloudinit_config: Remove length validation to allow empty content string in part blocks ([#105](https://github.com/hashicorp/terraform-provider-cloudinit/issues/105))

## 2.3.1 (February 22, 2023)

BUG FIXES:

* cloudinit_config: Fixed handling of unknown values in `part` blocks ([#103](https://github.com/hashicorp/terraform-provider-cloudinit/issues/103))

## 2.3.0 (February 22, 2023)

NOTES:

* provider: Rewritten to use the [`terraform-plugin-framework`](https://www.terraform.io/plugin/framework) ([#96](https://github.com/hashicorp/terraform-provider-cloudinit/issues/96))

## 2.2.0 (February 19, 2021)

Binary releases of this provider will now include the darwin-arm64 platform. This version contains no further changes.

## 2.1.0 (November 26, 2020)

NEW FEATURES:

* MIMEBOUNDARY can now be customised with `boundary` ([#7](https://github.com/hashicorp/terraform-provider-cloudinit/issues/7)).

## 2.0.0 (October 12, 2020)

Binary releases of this provider will now include the linux-arm64 platform.

BREAKING CHANGES:

* Upgrade to version 2 of the Terraform Plugin SDK, which drops support for Terraform 0.11. This provider will continue to work as expected for users of Terraform 0.11, which will not download the new version. ([#3](https://github.com/hashicorp/terraform-provider-cloudinit/issues/3))

## 1.0.0 (April 14, 2020)

Initial release. This provider exposes one resource, cloudinit_config, which is identical to the template_cloudinit_config resource in terraform-provider-template.
