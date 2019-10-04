package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/kmoe/terraform-provider-cloudinit/cloudinit"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cloudinit.Provider})
}
