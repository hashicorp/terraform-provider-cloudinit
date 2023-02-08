# cloud-init Provider Design

 Cloud-init is a commonly-used startup configuration utility for cloud compute instances. The cloud-init provider offers functionality to render a MIME multi-part file for use with cloud-init. Using a MIME multi-part file, the user can specify more than one type of data for cloud-init to consume. If you only have one type of data, you can leverage the built-in [`templatefile`](https://www.terraform.io/docs/configuration/functions/templatefile.html) function and a static file (like `.yml`).

Below we have a collection of _Goals_ and _Patterns_: they represent the guiding principles applied during the
development of this provider. Some are in place, others are ongoing processes, others are still just inspirational.

## Goals

* [_Stability over features_](.github/CONTRIBUTING.md)
* Provide a managed resource and data source to generate a cloud-init MIME multi-part file

## Patterns

Specific to this provider:

* The managed resource and data source use the same underlying code to generate the MIME multi-part file.

General to development:

* **Avoid repetition**: the entities managed can sometimes require similar pieces of logic and/or schema to be realised.
  When this happens it's important to keep the code shared in communal sections, so to avoid having to modify code in
  multiple places when they start changing.
* **Test expectations as well as bugs**: While it's typical to write tests to exercise a new functionality, it's key to
  also provide tests for issues that get identified and fixed, so to prove resolution as well as avoid regression.
* **Automate boring tasks**: Processes that are manual, repetitive and can be automated, should be. In addition to be a
  time-saving practice, this ensures consistency and reduces human error (ex. static code analysis).
* **Semantic versioning**: Adhering to HashiCorp's own
  [Versioning Specification](https://www.terraform.io/plugin/sdkv2/best-practices/versioning#versioning-specification)
  ensures we provide a consistent practitioner experience, and a clear process to deprecation and decommission.