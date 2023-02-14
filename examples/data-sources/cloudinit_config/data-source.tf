data "cloudinit_config" "foobar" {
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
