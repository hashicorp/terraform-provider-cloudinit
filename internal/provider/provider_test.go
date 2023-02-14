package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

var testProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
	"cloudinit": providerserver.NewProtocol5WithError(New()),
}
