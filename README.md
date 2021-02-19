Terraform Cloud-init Provider
==================

[Terraform](https://www.terraform.io) provider for rendering [cloud-init](https://cloudinit.readthedocs.io) configurations.

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

This provider is intended to replace the [template provider](https://www.terraform.io/docs/providers/template/). General templating can now be achieved through [the `templatefile` function](https://www.terraform.io/docs/configuration/functions/templatefile.html), without creating a separate data resource. 

The cloud-init Terraform provider exposes the `cloudinit_config` data source, previously available as the [`template_cloudinit_config` resource in the template provider](https://www.terraform.io/docs/providers/template/d/cloudinit_config.html).

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.16 (to build the provider plugin)


Using the provider
----------------------

The `cloudinit_config` data source renders a cloud-init config given in HCL form to the MIME-multipart form required by cloud-init.


Example configuration:
```
data "cloudinit_config" "foo" {
  gzip = false
  base64_encode = false

  part {
    content_type = "text/x-shellscript"
    content = "baz"
  }
}
```

This renders to:

```
Content-Type: multipart/mixed; boundary=\"MIMEBOUNDARY\"\nMIME-Version: 1.0\r\n\r\n--MIMEBOUNDARY\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/x-shellscript\r\nMime-Version: 1.0\r\n\r\nbaz\r\n--MIMEBOUNDARY--\r\n
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.16+ is *required*). You'll also need to set up a [GOPATH](http://golang.org/doc/code.html#GOPATH) correctly, as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-cloudinit
...
```

To test the provider, you can run `make test`.

```sh
$ make test
```

To run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources and often cost money to run.

```sh
$ make testacc
```
