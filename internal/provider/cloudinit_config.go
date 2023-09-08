// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	ID           types.String `tfsdk:"id"`
	Parts        types.List   `tfsdk:"part"` // configPartModel
	Gzip         types.Bool   `tfsdk:"gzip"`
	Base64Encode types.Bool   `tfsdk:"base64_encode"`
	Boundary     types.String `tfsdk:"boundary"`
	Rendered     types.String `tfsdk:"rendered"`
}

type configPartModel struct {
	ContentType types.String `tfsdk:"content_type"`
	Content     types.String `tfsdk:"content"`
	FileName    types.String `tfsdk:"filename"`
	MergeType   types.String `tfsdk:"merge_type"`
}

func (c *configModel) setDefaults(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	if c.Gzip.IsNull() {
		c.Gzip = types.BoolValue(true)
	}
	if c.Base64Encode.IsNull() {
		c.Base64Encode = types.BoolValue(true)
	}
	if c.Boundary.IsNull() {
		c.Boundary = types.StringValue("MIMEBOUNDARY")
	}

	if c.Parts.IsNull() || c.Parts.IsUnknown() {
		return diags
	}

	var configParts []configPartModel
	diags.Append(c.Parts.ElementsAs(ctx, &configParts, false)...)
	if diags.HasError() {
		return diags
	}

	for i, part := range configParts {
		if part.ContentType.IsNull() || part.ContentType.ValueString() == "" {
			configParts[i].ContentType = types.StringValue("text/plain")
		}
	}

	partsList, convertDiags := types.ListValueFrom(ctx, c.Parts.ElementType(ctx), configParts)
	diags.Append(convertDiags...)
	if diags.HasError() {
		return diags
	}

	c.Parts = partsList

	return diags
}

func (c configModel) validate(ctx context.Context) diag.Diagnostics {
	var diags diag.Diagnostics

	if c.Gzip.IsUnknown() || c.Base64Encode.IsUnknown() {
		return diags
	}
	diags.Append(c.setDefaults(ctx)...)
	if diags.HasError() {
		return diags
	}

	if c.Gzip.ValueBool() && !c.Base64Encode.ValueBool() {
		diags.AddAttributeError(
			path.Root("base64_encode"),
			"Invalid Attribute Configuration",
			"Expected base64_encode to be set to true when gzip is true.",
		)
	}

	return diags
}

func (c *configModel) update(ctx context.Context) diag.Diagnostics {
	var buffer bytes.Buffer
	var diags diag.Diagnostics
	var err error

	// cloudinit Provider 'v2.2.0' doesn't actually set default values in state properly, so we need to make sure
	// that we don't use any known empty values from previous versions of state
	diags.Append(c.setDefaults(ctx)...)
	if diags.HasError() {
		return diags
	}

	var configParts []configPartModel
	diags.Append(c.Parts.ElementsAs(ctx, &configParts, false)...)
	if diags.HasError() {
		return diags
	}

	if c.Gzip.ValueBool() {
		gzipWriter := gzip.NewWriter(&buffer)

		err = renderPartsToWriter(ctx, c.Boundary.ValueString(), configParts, gzipWriter)

		gzipWriter.Close()
	} else {
		err = renderPartsToWriter(ctx, c.Boundary.ValueString(), configParts, &buffer)
	}

	if err != nil {
		diags.AddError("Unable to render cloudinit config to MIME multi-part file", err.Error())
		return diags
	}

	output := ""
	if c.Base64Encode.ValueBool() {
		output = base64.StdEncoding.EncodeToString(buffer.Bytes())
	} else {
		output = buffer.String()
	}

	c.ID = types.StringValue(strconv.Itoa(hashcode.String(output)))
	c.Rendered = types.StringValue(output)

	return diags
}

func is7bit(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 0x7F {
			return false
		}
	}
	return true
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
		encode_to_base64 := !is7bit(part.Content.ValueString())
		var cte string
		if encode_to_base64 {
			cte = "base64"
		} else {
			cte = "7bit"
		}

		header := textproto.MIMEHeader{}

		header.Set("Content-Type", part.ContentType.ValueString())
		header.Set("MIME-Version", "1.0")
		header.Set("Content-Transfer-Encoding", cte)

		if part.FileName.ValueString() != "" {
			header.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, part.FileName.ValueString()))
		}

		if part.MergeType.ValueString() != "" {
			header.Set("X-Merge-Type", part.MergeType.ValueString())
		}

		var err error

		partWriter, err := mimeWriter.CreatePart(header)
		if err != nil {
			return err
		}

		if encode_to_base64 {
			base64Writer := base64.NewEncoder(base64.StdEncoding, partWriter)
			defer base64Writer.Close()
			_, err = base64Writer.Write([]byte(part.Content.ValueString()))
		} else {
			_, err = partWriter.Write([]byte(part.Content.ValueString()))
		}
		if err != nil {
			return err
		}
	}

	return nil
}
