package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hashicorp/terraform-provider-cloudinit/cloudinit"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cloudinit.Provider})
}
