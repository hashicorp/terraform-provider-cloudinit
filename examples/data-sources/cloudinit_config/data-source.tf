data "cloudinit_config" "foobar" {
  gzip          = false
  base64_encode = false

  part {
    content = "foo"

    content_type = "text/x-shellscript"
    filename     = "foo.sh"
  }

  part {
    content = "bar"

    content_type = "text/x-shellscript"
    filename     = "bar.sh"
  }
}
