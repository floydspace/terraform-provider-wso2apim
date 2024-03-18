---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "wso2apim_application_key_mapping Resource - wso2apim"
subcategory: ""
description: |-
  Manages a WSO2 API Manager Application Keys.
---

# wso2apim_application_key_mapping (Resource)

Manages a WSO2 API Manager Application Keys.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_id` (String) Application ID.
- `key_manager` (String) Key Manager.
- `key_type` (String) Application Key Type.
- `supported_grant_types` (List of String) Supported Grant Types.

### Optional

- `scopes` (List of String) Scopes.
- `validity_time` (Number) Validity Time.

### Read-Only

- `consumer_key` (String) Consumer Key.
- `consumer_secret` (String, Sensitive) Consumer Secret.
- `id` (String) Application Key Mapping ID.
- `key_state` (String) Application Key State.
- `last_updated` (String) Last updated timestamp.

## Import

Import is supported using the following syntax:

```shell
# Application can be imported by specifying the application id and key mapping id.
terraform import wso2apim_application_key_mapping.example 00000000-0000-0000-0000-000000000001,00000000-0000-0000-0000-000000000002
```