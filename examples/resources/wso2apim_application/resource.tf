# Manage example WSO2 API Manager Application
resource "wso2apim_application" "example" {
  name              = "foo-service"
  description       = "This is a foo service"
  throttling_policy = "Unlimited"
  token_type        = "JWT"
}
