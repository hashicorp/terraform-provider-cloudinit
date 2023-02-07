package provider

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
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
	"github.com/hashicorp/terraform-provider-cloudinit/internal/hashcode"
)

var (
	_ datasource.DataSourceWithValidateConfig = (*configDataSource)(nil)
)

type configDataSource struct{}

// TODO: better names.
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
			// TODO: is this path correct?
			path.Root("base64_encode"),
			"Invalid Attribute Configuration",
			"Expected base64_encode to be set to true when gzip is true.",
		)
	}
}

// TODO: Add descriptions and docs.
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
							Optional: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"content": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"filename": schema.StringAttribute{
							Optional: true,
						},
						"merge_type": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"gzip": schema.BoolAttribute{
				Optional: true,
			},
			"base64_encode": schema.BoolAttribute{
				Optional: true,
			},
			"boundary": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"rendered": schema.StringAttribute{
				Computed:    true,
				Description: "rendered cloudinit configuration",
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *configDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config configDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setDefaultValues(&config)

	renderedConfig, err := renderCloudinitConfig(&config)
	if err != nil {
		// TODO: add detail here
		resp.Diagnostics.AddError("error:", err.Error())
	}

	// TODO: should i map old struct to a new struct here?
	config.ID = types.StringValue(strconv.Itoa(hashcode.String(renderedConfig)))
	config.Rendered = types.StringValue(renderedConfig)

	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}

// NOTE: Currently it's not possible to specify default values
// against attributes of data sources in the schema.
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

// TODO: refactor
func renderCloudinitConfig(d *configDataSourceModel) (string, error) {
	var buffer bytes.Buffer

	var err error

	if d.Gzip.ValueBool() {
		gzipWriter := gzip.NewWriter(&buffer)
		err = renderPartsToWriter(d.Boundary.ValueString(), d.Parts, gzipWriter)
		err = gzipWriter.Close()
		if err != nil {
			return "", err
		}
	} else {
		err = renderPartsToWriter(d.Boundary.ValueString(), d.Parts, &buffer)
	}
	if err != nil {
		return "", err
	}

	output := ""
	if d.Base64Encode.ValueBool() {
		output = base64.StdEncoding.EncodeToString(buffer.Bytes())
	} else {
		output = buffer.String()
	}

	return output, nil
}

// TODO: refactor.
func renderPartsToWriter(mimeBoundary string, parts []configPartModel, writer io.Writer) error {
	mimeWriter := multipart.NewWriter(writer)
	defer func() {
		err := mimeWriter.Close()
		if err != nil {
			// TODO: switch to diag
			log.Printf("[WARN] Error closing mimewriter: %s", err)
		}
	}()

	// we need to set the boundary explicitly, otherwise the boundary is random
	// and this causes terraform to complain about the resource being different
	if err := mimeWriter.SetBoundary(mimeBoundary); err != nil {
		return err
	}

	writer.Write([]byte(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n", mimeWriter.Boundary())))
	writer.Write([]byte("MIME-Version: 1.0\r\n\r\n"))

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
