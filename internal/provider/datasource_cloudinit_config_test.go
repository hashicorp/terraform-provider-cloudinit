package provider

import (
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testProviders = map[string]*schema.Provider{
	"cloudinit": New(),
}

func TestRender(t *testing.T) {
	testCases := []struct {
		Name          string
		ResourceBlock string
		Expected      string
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
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			r.UnitTest(t, r.TestCase{
				Providers: testProviders,
				Steps: []r.TestStep{
					{
						Config: tt.ResourceBlock,
						Check: r.ComposeTestCheckFunc(
							r.TestCheckResourceAttr("data.cloudinit_config.foo", "rendered", tt.Expected),
						),
					},
				},
			})
		})
	}
}

// From GH-13572, Correctly handle panic on a misconfigured cloudinit part.
func TestRender_handlePanic(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config:      testCloudInitConfig_misconfiguredParts,
				ExpectError: regexp.MustCompile("Unable to parse parts in cloudinit resource declaration"),
			},
		},
	})
}

var testCloudInitConfig_misconfiguredParts = `
data "cloudinit_config" "foo" {
  part {
    content = ""
  }
}
`
