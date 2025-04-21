## 2.3.7 (April 21, 2025)

NOTES:

* Update dependencies ([#339](https://github.com/hashicorp/terraform-provider-cloudinit/issues/339))

## 2.3.7-alpha1 (March 04, 2025)

NOTES:

* all: This release is being used to test new build and release actions.

## 2.3.6 (February 27, 2025)

NOTES:

* all: This release contains no functionality changes. It is being used to fix release metadata in the registry from the previous alpha1 release. ([#318](https://github.com/hashicorp/terraform-provider-cloudinit/issues/318))

## 2.3.6-alpha1 (December 05, 2024)

NOTES:

* all: This release contains no functionality changes. It is released using new build and release Actions. ([#293](https://github.com/hashicorp/terraform-provider-cloudinit/issues/293))

## 2.3.5 (September 10, 2024)

NOTES:

* all: This release introduces no functional changes. It does however include dependency updates which address upstream CVEs. ([#263](https://github.com/hashicorp/terraform-provider-cloudinit/issues/263))

## 2.3.4 (April 22, 2024)

NOTES:

* all: This release contains no functionality changes, only the inclusion of the LICENSE file in the release archives ([#228](https://github.com/hashicorp/terraform-provider-cloudinit/issues/228))

## 2.3.3 (November 29, 2023)

NOTES:

* This release introduces no functional changes. It does however include dependency updates which address upstream CVEs. ([#186](https://github.com/hashicorp/terraform-provider-cloudinit/issues/186))

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
