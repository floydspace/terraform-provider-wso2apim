package wso2apim

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "wso2apim_application" "test" {
	name              = "foo-app1"
	throttling_policy = "Unlimited"
	description       = "This is a foo application"

	attributes = {
		"serviceNowManagedByGroup" = "test"
		"serviceNowSupportGroup"   = "test"
	}
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wso2apim_application.test", "name", "foo-app1"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "description", "This is a foo application"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "throttling_policy", "Unlimited"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "token_type", "JWT"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "attributes.serviceNowManagedByGroup", "test"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "attributes.serviceNowSupportGroup", "test"),
					resource.TestCheckResourceAttrSet("wso2apim_application.test", "id"),
					resource.TestCheckResourceAttrSet("wso2apim_application.test", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "wso2apim_application.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "wso2apim_application" "test" {
	name              = "foo-app1"
	throttling_policy = "Unlimited"
	description       = "This is a foo application ?"

	attributes = {
		"serviceNowManagedByGroup" = "test"
		"serviceNowSupportGroup"   = "test"
	}
}
						`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wso2apim_application.test", "name", "foo-app1"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "description", "This is a foo application ?"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "throttling_policy", "Unlimited"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "token_type", "JWT"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "attributes.serviceNowManagedByGroup", "test"),
					resource.TestCheckResourceAttr("wso2apim_application.test", "attributes.serviceNowSupportGroup", "test"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
