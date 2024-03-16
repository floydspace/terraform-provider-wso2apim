# Manage example WSO2 API Manager Application Key Mapping
resource "wso2apim_application_key_mapping" "example" {
  application_id = wso2apim_application.example.id
  key_manager    = data.wso2apim_key_manager.example.id
  key_type       = "PRODUCTION"
  supported_grant_types = [
    "client_credentials",
    "password",
    "refresh_token"
  ]
  scopes        = ["default"]
  validity_time = 3600
}
