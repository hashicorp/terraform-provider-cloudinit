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

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-cloudinit/internal/hashcode"
)

// Model and functionality of data source and resource are equivalent.
type configModel struct {
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

func validateConfigModel(config configModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if config.Gzip.IsUnknown() || config.Base64Encode.IsUnknown() {
		return diags
	}
	setDefaultValues(&config)

	if config.Gzip.ValueBool() && !config.Base64Encode.ValueBool() {
		diags.AddAttributeError(
			path.Root("base64_encode"),
			"Invalid Attribute Configuration",
			"Expected base64_encode to be set to true when gzip is true.",
		)
	}

	return diags
}

func updateConfigModel(ctx context.Context, config *configModel) diag.Diagnostics {
	var diags diag.Diagnostics

	renderedConfig, err := renderCloudinitConfig(ctx, config)
	if err != nil {
		diags.AddError("Unable to render cloudinit config", err.Error())
		return diags
	}

	config.ID = types.StringValue(strconv.Itoa(hashcode.String(renderedConfig)))
	config.Rendered = types.StringValue(renderedConfig)

	return diags
}

func renderCloudinitConfig(ctx context.Context, c *configModel) (string, error) {
	var buffer bytes.Buffer
	var err error

	if c.Gzip.ValueBool() {
		gzipWriter := gzip.NewWriter(&buffer)
		err = renderPartsToWriter(ctx, c.Boundary.ValueString(), c.Parts, gzipWriter)

		gzipWriter.Close()
	} else {
		err = renderPartsToWriter(ctx, c.Boundary.ValueString(), c.Parts, &buffer)
	}

	if err != nil {
		return "", fmt.Errorf("error writing part block to MIME multi-part file: %w", err)
	}

	output := ""
	if c.Base64Encode.ValueBool() {
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

func setDefaultValues(c *configModel) {
	if c.Gzip.IsNull() {
		c.Gzip = types.BoolValue(true)
	}
	if c.Base64Encode.IsNull() {
		c.Base64Encode = types.BoolValue(true)
	}
	if c.Boundary.IsNull() {
		c.Boundary = types.StringValue("MIMEBOUNDARY")
	}

	for i, part := range c.Parts {
		if part.ContentType.IsNull() {
			c.Parts[i].ContentType = types.StringValue("text/plain")
		}
	}
}
