package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-cloudinit/internal/provider/attribute_plan_modifier_bool"
	"github.com/hashicorp/terraform-provider-cloudinit/internal/provider/attribute_plan_modifier_string"
)

var (
	_ resource.ResourceWithValidateConfig = (*configResource)(nil)
)

type configResource struct{}

func (r *configResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (r *configResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config configModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateConfigModel(config)...)
}

func (r *configResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
								attribute_plan_modifier_string.DefaultValue(types.StringValue("text/plain")),
							},
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "A MIME-style content type to report in the header for the part. Defaults to `text/plain`",
						},
						"content": schema.StringAttribute{
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Required:            true,
							MarkdownDescription: "Body content for the part.",
						},
						"filename": schema.StringAttribute{
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Optional:            true,
							MarkdownDescription: "A filename to report in the header for the part.",
						},
						"merge_type": schema.StringAttribute{
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
							Optional: true,
							MarkdownDescription: "A value for the `X-Merge-Type` header of the part, to control " +
								"[cloud-init merging behavior](https://cloudinit.readthedocs.io/en/latest/reference/merging.html).",
						},
					},
				},
				// TODO: add note about this being required in docs?
				MarkdownDescription: "A nested block type which adds a file to the generated cloud-init configuration. Use multiple " +
					"`part` blocks to specify multiple files, which will be included in order of declaration in the final MIME document.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"gzip": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					attribute_plan_modifier_bool.DefaultValue(types.BoolValue(true)),
				},
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specify whether or not to gzip the `rendered` output. Defaults to `true`.",
			},
			"base64_encode": schema.BoolAttribute{
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					attribute_plan_modifier_bool.DefaultValue(types.BoolValue(true)),
				},
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specify whether or not to base64 encode the `rendered` output. Defaults to `true`, and cannot be disabled if gzip is `true`.",
			},
			"boundary": schema.StringAttribute{
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					attribute_plan_modifier_string.DefaultValue(types.StringValue("MIMEBOUNDARY")),
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
		MarkdownDescription: `**NOTE**: This resource is deprecated, use data source instead.`,
		DeprecationMessage:  `**NOTE**: This resource is deprecated, use data source instead.`,
	}
}

func (r *configResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var newState configModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateConfigModel(ctx, &newState)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *configResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var newState configModel

	resp.Diagnostics.Append(req.State.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(updateConfigModel(ctx, &newState)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *configResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *configResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
