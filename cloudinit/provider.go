package cloudinit

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
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
