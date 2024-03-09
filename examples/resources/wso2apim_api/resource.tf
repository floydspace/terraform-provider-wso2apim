# Manage example WSO2 API Manager Api
resource "wso2apim_api" "example" {
  name        = "foo-api"
  description = "This is a foo API"
  context     = "/foo"
  version     = "v1"
}
