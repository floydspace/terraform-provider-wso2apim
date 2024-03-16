# Manage example WSO2 API Manager Subscription
resource "wso2apim_subscription" "example" {
  application_id    = wso2apim_application.example.id
  api_id            = wso2apim_api.example.id
  throttling_policy = "Unlimited"
}
