package wso2apim

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApiResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "wso2apim_api" "test" {
	name        = "foo-api2"
	description = "This is a foo API"
	context     = "/bar2"
	version     = "v1"
	operations = [{
		target = "/graphql"
		verb   = "POST"
	}]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wso2apim_api.test", "name", "foo-api2"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "description", "This is a foo API"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "context", "/bar2"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "version", "v1"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "operations.0.target", "/graphql"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "operations.0.verb", "POST"),
					resource.TestCheckResourceAttrSet("wso2apim_api.test", "id"),
					resource.TestCheckResourceAttrSet("wso2apim_api.test", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "wso2apim_api.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "wso2apim_api" "test" {
	name        = "foo-api2"
	description = "This is a foo API ?"
	context     = "/bar2"
	version     = "v1"
	operations = [{
		target = "/graphql"
		verb   = "POST"
	}]
}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wso2apim_api.test", "name", "foo-api2"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "description", "This is a foo API ?"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "context", "/bar2"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "version", "v1"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "operations.0.target", "/graphql"),
					resource.TestCheckResourceAttr("wso2apim_api.test", "operations.0.verb", "POST"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
