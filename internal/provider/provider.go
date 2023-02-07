package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = (*cloudinitProvider)(nil)
)

type cloudinitProvider struct{}

func New() provider.Provider {
	return &cloudinitProvider{}
}

func (p *cloudinitProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudinit"
}

// TODO: Add descriptions and docs
func (p *cloudinitProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *cloudinitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *cloudinitProvider) Resources(ctx context.Context) []func() resource.Resource {
	// TODO: Resource shim was in old code?
	// 		ResourcesMap: map[string]*schema.Resource{
	// 			"cloudinit_config": schema.DataSourceResourceShim(
	// 				"cloudinit_config",
	// 				dataSourceCloudinitConfig(),
	// 			),
	// 		},
	return []func() resource.Resource{}
}

func (p *cloudinitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &configDataSource{}
		},
	}
}
