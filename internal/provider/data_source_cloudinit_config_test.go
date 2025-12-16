// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConfigDataSourceRender(t *testing.T) {
	testCases := []struct {
		Name            string
		DataSourceBlock string
		Expected        string
	}{
		{
			"no gzip or b64 - basic content",
			`data "cloudinit_config" "foo" {
				gzip = false
				base64_encode = false

				part {
					content_type = "text/x-shellscript"
					content = "baz"
				}
			}`,
			"Content-Type: multipart/mixed; boundary=\"MIMEBOUNDARY\"\nMIME-Version: 1.0\r\n\r\n--MIMEBOUNDARY\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/x-shellscript\r\nMime-Version: 1.0\r\n\r\nbaz\r\n--MIMEBOUNDARY--\r\n",
		},
		{
			"no gzip or b64 - basic content - default to text plain",
			`data "cloudinit_config" "foo" {
				gzip = false
				base64_encode = false

				part {
					content = "baz"
				}
			}`,
			"Content-Type: multipart/mixed; boundary=\"MIMEBOUNDARY\"\nMIME-Version: 1.0\r\n\r\n--MIMEBOUNDARY\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/plain\r\nMime-Version: 1.0\r\n\r\nbaz\r\n--MIMEBOUNDARY--\r\n",
		},
		{
			"no gzip or b64 - content with filename",
			`data "cloudinit_config" "foo" {
				gzip = false
				base64_encode = false

				part {
					content_type = "text/x-shellscript"
					content = "baz"
					filename = "foobar.sh"
				}
			}`,
			"Content-Type: multipart/mixed; boundary=\"MIMEBOUNDARY\"\nMIME-Version: 1.0\r\n\r\n--MIMEBOUNDARY\r\nContent-Disposition: attachment; filename=\"foobar.sh\"\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/x-shellscript\r\nMime-Version: 1.0\r\n\r\nbaz\r\n--MIMEBOUNDARY--\r\n",
		},
		{
			"no gzip or b64 - two parts, basic content",
			`data "cloudinit_config" "foo" {
				gzip = false
				base64_encode = false

				part {
					content_type = "text/x-shellscript"
					content = "baz"
				}
				part {
					content_type = "text/x-shellscript"
					content = "ffbaz"
				}
			}`,
			"Content-Type: multipart/mixed; boundary=\"MIMEBOUNDARY\"\nMIME-Version: 1.0\r\n\r\n--MIMEBOUNDARY\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/x-shellscript\r\nMime-Version: 1.0\r\n\r\nbaz\r\n--MIMEBOUNDARY\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/x-shellscript\r\nMime-Version: 1.0\r\n\r\nffbaz\r\n--MIMEBOUNDARY--\r\n",
		},
		{
			"no gzip or b64 - with boundary separator",
			`data "cloudinit_config" "foo" {
				gzip = false
				base64_encode = false
				boundary = "//"

				part {
					content_type = "text/x-shellscript"
					content = "baz"
				}
			}`,
			"Content-Type: multipart/mixed; boundary=\"//\"\nMIME-Version: 1.0\r\n\r\n--//\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/x-shellscript\r\nMime-Version: 1.0\r\n\r\nbaz\r\n--//--\r\n",
		},
		{
			"no gzip or b64 - two parts - all fields",
			`data "cloudinit_config" "foo" {
				gzip = false
				base64_encode = false

				part {
					content_type = "text/x-shellscript"
					content = "foo1"
					filename = "foofile1.txt"
					merge_type = "list()+dict()+str()"
				}

				part {
					content_type = "text/x-shellscript"
					content = "bar1"
					filename = "barfile1.txt"
					merge_type = "list()+dict()+str()"
				}
			}`,
			"Content-Type: multipart/mixed; boundary=\"MIMEBOUNDARY\"\nMIME-Version: 1.0\r\n\r\n--MIMEBOUNDARY\r\nContent-Disposition: attachment; filename=\"foofile1.txt\"\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/x-shellscript\r\nMime-Version: 1.0\r\nX-Merge-Type: list()+dict()+str()\r\n\r\nfoo1\r\n--MIMEBOUNDARY\r\nContent-Disposition: attachment; filename=\"barfile1.txt\"\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/x-shellscript\r\nMime-Version: 1.0\r\nX-Merge-Type: list()+dict()+str()\r\n\r\nbar1\r\n--MIMEBOUNDARY--\r\n",
		},
		{
			"no gzip - b64 encoded - basic content",
			`data "cloudinit_config" "foo" {
				gzip = false
				base64_encode = true

				part {
					content_type = "text/x-shellscript"
					content = "heythere"
				}
			}`,
			"Q29udGVudC1UeXBlOiBtdWx0aXBhcnQvbWl4ZWQ7IGJvdW5kYXJ5PSJNSU1FQk9VTkRBUlkiCk1JTUUtVmVyc2lvbjogMS4wDQoNCi0tTUlNRUJPVU5EQVJZDQpDb250ZW50LVRyYW5zZmVyLUVuY29kaW5nOiA3Yml0DQpDb250ZW50LVR5cGU6IHRleHQveC1zaGVsbHNjcmlwdA0KTWltZS1WZXJzaW9uOiAxLjANCg0KaGV5dGhlcmUNCi0tTUlNRUJPVU5EQVJZLS0NCg==",
		},
		{
			"gzip compression - basic content",
			`data "cloudinit_config" "foo" {
				gzip = true

				part {
					content_type = "text/x-shellscript"
					content = "heythere"
				}
			}`,
			"H4sIAAAAAAAA/2TNuwrCQBCF4X5h32FJP0YrIWLhJYVFFEQFy1xGM5DMhtkJJG8vWkjQ8sDP+XaeFVnhMnaYuLZvlLpcNG5pwGrlCt9zlcu4jrJDlm5P1+N+c75H5r3ghhLIc+IWs7k11gBMI2u+35JzeKBAyqWviJ+JWxakk+CDKw4aDxBqbJpQCnVqTUYt/jk1jlqj4K8IYM0rAAD//0u6BO3QAAAA",
		},
		// Initial fix of panic when part block wasn't parsable
		// https://github.com/hashicorp/terraform/issues/13572
		//
		// Move to framework fixed parsing error, but introduced incorrect validation,  empty content is valid
		// https://github.com/hashicorp/terraform-provider-cloudinit/issues/104
		{
			"empty content field in part block",
			`data "cloudinit_config" "foo" {
				part {
					content = ""
				}
			}`,
			"H4sIAAAAAAAA/3LOzytJzSvRDaksSLVSyC3NKcksSCwq0c/NrEhNsVZIyi/NS0ksqrRV8vX0dXXyD/VzcQyKVOIC8XTDUouKM/PzrBQM9Qx4uXi5dHWRFfFywc0uSswrTkst0nXNS85PycxLt1IwT8osQVIAtrwktaJEvyAnMTOPl8s3MzcVw3x0G3R1ebkAAQAA///dakG4wAAAAA==",
		},
		{
			"empty content field in one part block, multiple part blocks",
			`data "cloudinit_config" "foo" {
				gzip = false
				base64_encode = false

				part {
					content_type = "text/cloud-config"
					content = "abc"
				}

				part {
					content = ""
				}

				part {
					content_type = "text/cloud-config"
					content = ""
				}
			}`,
			"Content-Type: multipart/mixed; boundary=\"MIMEBOUNDARY\"\nMIME-Version: 1.0\r\n\r\n--MIMEBOUNDARY\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/cloud-config\r\nMime-Version: 1.0\r\n\r\nabc\r\n--MIMEBOUNDARY\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/plain\r\nMime-Version: 1.0\r\n\r\n\r\n--MIMEBOUNDARY\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/cloud-config\r\nMime-Version: 1.0\r\n\r\n\r\n--MIMEBOUNDARY--\r\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			r.UnitTest(t, r.TestCase{
				ProtoV5ProviderFactories: testProtoV5ProviderFactories,
				Steps: []r.TestStep{
					{
						Config: tt.DataSourceBlock,
						Check: r.ComposeTestCheckFunc(
							r.TestCheckResourceAttr("data.cloudinit_config.foo", "rendered", tt.Expected),
						),
					},
				},
			})
		})
	}
}

func TestConfigDataSourceRender_handleErrors(t *testing.T) {
	testCases := []struct {
		Name            string
		DataSourceBlock string
		ErrorMatch      *regexp.Regexp
	}{
		{
			"base64 can't be false when gzip is true",
			`data "cloudinit_config" "foo" {
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
			`data "cloudinit_config" "foo" {
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
						Config:      tt.DataSourceBlock,
						ExpectError: tt.ErrorMatch,
					},
				},
			})
		})
	}
}

func TestConfigDataSource_UpgradeFromVersion2_2_0(t *testing.T) {
	testCases := []struct {
		Name            string
		DataSourceBlock string
		Expected        string
	}{
		{
			"multiple parts in cloudinit config",
			`data "cloudinit_config" "foo" {
				gzip = false
				base64_encode = false

				part {
					content_type = "text/x-shellscript"
					content = "foo1"
					filename = "foofile1.txt"
					merge_type = "list()+dict()+str()"
				}

				part {
					content = "bar1"
					filename = "barfile1.txt"
					merge_type = "list()+dict()+str()"
				}
			}`,
			"Content-Type: multipart/mixed; boundary=\"MIMEBOUNDARY\"\nMIME-Version: 1.0\r\n\r\n--MIMEBOUNDARY\r\nContent-Disposition: attachment; filename=\"foofile1.txt\"\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/x-shellscript\r\nMime-Version: 1.0\r\nX-Merge-Type: list()+dict()+str()\r\n\r\nfoo1\r\n--MIMEBOUNDARY\r\nContent-Disposition: attachment; filename=\"barfile1.txt\"\r\nContent-Transfer-Encoding: 7bit\r\nContent-Type: text/plain\r\nMime-Version: 1.0\r\nX-Merge-Type: list()+dict()+str()\r\n\r\nbar1\r\n--MIMEBOUNDARY--\r\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			r.UnitTest(t, r.TestCase{
				Steps: []r.TestStep{
					{
						ExternalProviders: map[string]r.ExternalProvider{
							"cloudinit": {
								VersionConstraint: "2.2.0",
								Source:            "hashicorp/cloudinit",
							},
						},
						Config: tt.DataSourceBlock,
						Check: r.ComposeTestCheckFunc(
							r.TestCheckResourceAttr("data.cloudinit_config.foo", "rendered", tt.Expected),
						),
					},
					{
						ProtoV5ProviderFactories: testProtoV5ProviderFactories,
						Config:                   tt.DataSourceBlock,
						Check: r.ComposeTestCheckFunc(
							r.TestCheckResourceAttr("data.cloudinit_config.foo", "rendered", tt.Expected),
						),
					},
				},
			})
		})
	}
}

// This test ensures that unknown values are being handled properly in the `part` block
// https://github.com/hashicorp/terraform-provider-cloudinit/issues/102
func TestConfigDataSource_HandleUnknown(t *testing.T) {
	testCases := []struct {
		Name            string
		DataSourceBlock string
		Expected        string
	}{
		{
			"issue 102 - handle unknown in validate caused by variable",
			`
			variable "config_types" {
				default = ["text/cloud-config", "text/cloud-config"]
			}
			  
			data "cloudinit_config" "foo" {
				gzip          = true
				base64_encode = true
				
				dynamic "part" {
				  for_each = var.config_types
				  content {
					content_type = part.value
					content      = <<-EOT
					#cloud-config
					test: ${part.value}
					EOT
				  }
				}
			}
			`,
			"H4sIAAAAAAAA/3LOzytJzSvRDaksSLVSyC3NKcksSCwq0c/NrEhNsVZIyi/NS0ksqrRV8vX0dXXyD/VzcQyKVOIC8XTDUouKM/PzrBQM9Qx4uXi5dHWRFfFywc0uSswrTkst0nXNS85PycxLt1IwT8osQVIAtrwktaJEPzknvzRFNzk/Ly0znZfLNzM3FcMaZWQ1XCWpxSVY9A525+jq8nIBAgAA//821u5zfAEAAA==",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			r.UnitTest(t, r.TestCase{
				Steps: []r.TestStep{
					{
						ProtoV5ProviderFactories: testProtoV5ProviderFactories,
						Config:                   tt.DataSourceBlock,
						Check: r.ComposeTestCheckFunc(
							r.TestCheckResourceAttr("data.cloudinit_config.foo", "rendered", tt.Expected),
						),
					},
				},
			})
		})
	}
}
