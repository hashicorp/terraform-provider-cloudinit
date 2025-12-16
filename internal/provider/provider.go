// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

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

func (p *cloudinitProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The cloud-init Terraform provider exposes the `cloudinit_config` data source, previously available as the " +
			"`template_cloudinit_config` resource [in the template provider](https://registry.terraform.io/providers/hashicorp/template/latest/docs/data-sources/cloudinit_config), " +
			"which renders a [multipart MIME configuration](https://cloudinit.readthedocs.io/en/latest/explanation/format.html#mime-multi-part-archive) " +
			"for use with [cloud-init](https://cloudinit.readthedocs.io/en/latest/).",
	}
}

func (p *cloudinitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *cloudinitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &configResource{}
		},
	}
}

func (p *cloudinitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &configDataSource{}
		},
	}
}
