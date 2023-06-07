// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build generate

package tools

import (
	// document generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	// copywrite header generation
	_ "github.com/hashicorp/copywrite"
)

//go:generate go run github.com/hashicorp/copywrite headers -d .. --config ../.copywrite.hcl
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir ..
