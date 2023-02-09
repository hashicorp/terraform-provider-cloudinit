resource "cloudinit_config" "foo" {
  gzip          = false
  base64_encode = false

  part {
    content = "baz"

    content_type = "text/x-shellscript"
    filename     = "foobar.sh"
  }
}
