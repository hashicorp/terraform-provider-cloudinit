// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConfigEphemeralResourceBasic(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: testProtoV5ProviderFactories,
		Steps: []r.TestStep{
			{
				Config: `
					ephemeral "cloudinit_config" "test" {
						gzip = false
						base64_encode = false

						part {
							content_type = "text/x-shellscript"
							content = "#!/bin/bash\necho 'hello world'"
						}
					}
				`,
			},
		},
	})
}

func TestConfigEphemeralResourceRender_handleErrors(t *testing.T) {
	testCases := []struct {
		Name          string
		ResourceBlock string
		ErrorMatch    *regexp.Regexp
	}{
		{
			"base64 can't be false when gzip is true",
			`ephemeral "cloudinit_config" "foo" {
				gzip = true
				base64_encode = false

				part {
				  content = "abc"
				}
			}`,
			regexp.MustCompile("Expected base64_encode to be set to true when gzip is true"),
		},
		{
			"at least one part is required",
			`ephemeral "cloudinit_config" "foo" {
				gzip = false
				base64_encode = false
			}`,
			regexp.MustCompile("part must have a configuration value"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			r.UnitTest(t, r.TestCase{
				ProtoV5ProviderFactories: testProtoV5ProviderFactories,
				Steps: []r.TestStep{
					{
						Config:      tt.ResourceBlock,
						ExpectError: tt.ErrorMatch,
					},
				},
			})
		})
	}
}

