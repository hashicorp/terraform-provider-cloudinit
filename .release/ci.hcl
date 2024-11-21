# Reference: https://github.com/hashicorp/crt-core-helloworld/blob/main/.release/ci.hcl (private repository)

schema = "2"

project "terraform-provider-cloudinit" {
  // team is currently unused and has no meaning
  // but is required to be non-empty by CRT orchestator
  team = "_UNUSED_"

  slack {
    notification_channel = "C02M018DV27" // #feed-tf-devex
  }

  github {
    organization     = "hashicorp"
    repository       = "terraform-provider-cloudinit"
    release_branches = ["main"]
  }
}

event "build" {
  action "build" {
    organization = "hashicorp"
    repository   = "terraform-provider-cloudinit"
    workflow     = "build"
  }
}

event "prepare" {
  # `prepare` is the Common Release Tooling (CRT) artifact processing workflow.
  # It prepares artifacts for potential promotion to staging and production.
  # For example, it scans and signs artifacts.

  depends = ["build"]

  action "prepare" {
    organization = "hashicorp"
    repository   = "crt-workflows-common"
    workflow     = "prepare"
    depends      = ["build"]
  }

  notification {
    on = "fail"
  }
}
