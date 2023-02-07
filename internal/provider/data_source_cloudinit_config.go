package provider

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-cloudinit/internal/hashcode"
)

var (
	_ datasource.DataSourceWithValidateConfig = (*configDataSource)(nil)
)

type configDataSource struct{}

type configDataSourceModel struct {
	ID           types.String      `tfsdk:"id"`
	Parts        []configPartModel `tfsdk:"part"`
	Gzip         types.Bool        `tfsdk:"gzip"`
	Base64Encode types.Bool        `tfsdk:"base64_encode"`
	Boundary     types.String      `tfsdk:"boundary"`
	Rendered     types.String      `tfsdk:"rendered"`
}

type configPartModel struct {
	ContentType types.String `tfsdk:"content_type"`
	Content     types.String `tfsdk:"content"`
	FileName    types.String `tfsdk:"filename"`
	MergeType   types.String `tfsdk:"merge_type"`
}

func (d *configDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (*configDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var config configDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Gzip.IsUnknown() || config.Base64Encode.IsUnknown() {
		return
	}
	setDefaultValues(&config)

	if config.Gzip.ValueBool() && !config.Base64Encode.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base64_encode"),
			"Invalid Attribute Configuration",
			"Expected base64_encode to be set to true when gzip is true.",
		)
	}
}

func (d *configDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"part": schema.ListNestedBlock{
				Validators: []validator.List{
					// TODO: make sure I switch the go.mod back to newly released validators after merging validators
					listvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"content_type": schema.StringAttribute{
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
							Optional:            true,
							MarkdownDescription: "A MIME-style content type to report in the header for the part. Defaults to `text/plain`",
						},
						"content": schema.StringAttribute{
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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
				// TODO: add note about this being required? what will show up?
				MarkdownDescription: "A nested block type which adds a file to the generated cloud-init configuration. Use multiple " +
					"`part` blocks to specify multiple files, which will be included in order of declaration in the final MIME document.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"gzip": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Specify whether or not to gzip the `rendered` output. Defaults to `true`.",
			},
			"base64_encode": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Specify whether or not to base64 encode the `rendered` output. Defaults to `true`, and cannot be disabled if gzip is `true`.",
			},
			"boundary": schema.StringAttribute{
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Optional:            true,
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
	var newState configDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setDefaultValues(&newState)

	renderedConfig, err := renderCloudinitConfig(ctx, &newState)
	if err != nil {
		resp.Diagnostics.AddError("Unable to render cloudinit config", err.Error())
		return
	}

	newState.ID = types.StringValue(strconv.Itoa(hashcode.String(renderedConfig)))
	newState.Rendered = types.StringValue(renderedConfig)

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

// NOTE: Currently it's not possible to specify default values against attributes of data sources in the schema.
func setDefaultValues(d *configDataSourceModel) {
	if d.Gzip.IsNull() {
		d.Gzip = types.BoolValue(true)
	}
	if d.Base64Encode.IsNull() {
		d.Base64Encode = types.BoolValue(true)
	}
	if d.Boundary.IsNull() {
		d.Boundary = types.StringValue("MIMEBOUNDARY")
	}

	for i, part := range d.Parts {
		if part.ContentType.IsNull() {
			d.Parts[i].ContentType = types.StringValue("text/plain")
		}
	}
}

func renderCloudinitConfig(ctx context.Context, d *configDataSourceModel) (string, error) {
	var buffer bytes.Buffer
	var err error

	if d.Gzip.ValueBool() {
		gzipWriter := gzip.NewWriter(&buffer)
		err = renderPartsToWriter(ctx, d.Boundary.ValueString(), d.Parts, gzipWriter)

		gzipWriter.Close()
	} else {
		err = renderPartsToWriter(ctx, d.Boundary.ValueString(), d.Parts, &buffer)
	}

	if err != nil {
		return "", fmt.Errorf("error writing part block to MIME multi-part file: %w", err)
	}

	output := ""
	if d.Base64Encode.ValueBool() {
		output = base64.StdEncoding.EncodeToString(buffer.Bytes())
	} else {
		output = buffer.String()
	}

	return output, nil
}

func renderPartsToWriter(ctx context.Context, mimeBoundary string, parts []configPartModel, writer io.Writer) error {
	mimeWriter := multipart.NewWriter(writer)
	defer func() {
		err := mimeWriter.Close()
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("error closing mimeWriter: %s", err))
		}
	}()

	// we need to set the boundary explicitly, otherwise the boundary is random
	// and this causes terraform to complain about the resource being different
	if err := mimeWriter.SetBoundary(mimeBoundary); err != nil {
		return err
	}

	_, err := writer.Write([]byte(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n", mimeWriter.Boundary())))
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte("MIME-Version: 1.0\r\n\r\n"))
	if err != nil {
		return err
	}

	for _, part := range parts {
		header := textproto.MIMEHeader{}

		header.Set("Content-Type", part.ContentType.ValueString())
		header.Set("MIME-Version", "1.0")
		header.Set("Content-Transfer-Encoding", "7bit")

		if part.FileName.ValueString() != "" {
			header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, part.FileName.ValueString()))
		}

		if part.MergeType.ValueString() != "" {
			header.Set("X-Merge-Type", part.MergeType.ValueString())
		}

		partWriter, err := mimeWriter.CreatePart(header)
		if err != nil {
			return err
		}

		_, err = partWriter.Write([]byte(part.Content.ValueString()))
		if err != nil {
			return err
		}
	}

	return nil
}
