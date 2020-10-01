package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"cloudinit_config": dataSourceCloudinitConfig(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudinit_config": schema.DataSourceResourceShim(
				"cloudinit_config",
				dataSourceCloudinitConfig(),
			),
		},
	}
}
