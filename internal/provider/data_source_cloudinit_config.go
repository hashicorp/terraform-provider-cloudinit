// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ datasource.DataSourceWithValidateConfig = (*configDataSource)(nil)
)

type configDataSource struct{}

func (d *configDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (d *configDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var cloudinitConfig configModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &cloudinitConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(cloudinitConfig.validate(ctx)...)
}

func (d *configDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"part": schema.ListNestedBlock{
				Validators: []validator.List{
					listvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"content_type": schema.StringAttribute{
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "A MIME-style content type to report in the header for the part. Defaults to `text/plain`",
						},
						"content": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Body content for the part.",
						},
						"filename": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "A filename to report in the header for the part.",
						},
						"merge_type": schema.StringAttribute{
							Optional: true,
							MarkdownDescription: "A value for the `X-Merge-Type` header of the part, to control " +
								"[cloud-init merging behavior](https://cloudinit.readthedocs.io/en/latest/reference/merging.html).",
						},
					},
				},
				MarkdownDescription: "A nested block type which adds a file to the generated cloud-init configuration. Use multiple " +
					"`part` blocks to specify multiple files, which will be included in order of declaration in the final MIME document.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"gzip": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specify whether or not to gzip the `rendered` output. Defaults to `true`.",
			},
			"base64_encode": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specify whether or not to base64 encode the `rendered` output. Defaults to `true`, and cannot be disabled if gzip is `true`.",
			},
			"boundary": schema.StringAttribute{
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specify the Writer's default boundary separator. Defaults to `MIMEBOUNDARY`.",
			},
			"rendered": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The final rendered multi-part cloud-init config.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "[CRC-32](https://pkg.go.dev/hash/crc32) checksum of `rendered` cloud-init config.",
			},
		},
		MarkdownDescription: "Renders a [multi-part MIME configuration](https://cloudinit.readthedocs.io/en/latest/explanation/format.html#mime-multi-part-archive) " +
			"for use with [cloud-init](https://cloudinit.readthedocs.io/en/latest/).\n\n" +
			"Cloud-init is a commonly-used startup configuration utility for cloud compute instances. It accepts configuration via provider-specific " +
			"user data mechanisms, such as `user_data` for Amazon EC2 instances. Multi-part MIME is one of the data formats it accepts. For more information, " +
			"see [User-Data Formats](https://cloudinit.readthedocs.io/en/latest/explanation/format.html) in the cloud-init manual.\n\n" +
			"This is not a generalized utility for producing multi-part MIME messages. It's feature set is specialized for cloud-init multi-part MIME messages.",
	}
}

func (d *configDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cloudinitConfig configModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &cloudinitConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(cloudinitConfig.update(ctx)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, cloudinitConfig)...)
}
