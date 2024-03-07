# terraform-provider-wso2apim

This is a terraform provider for managing APIs on WSO2 API Manager.

## Example

```tf
provider "wso2apim" {
  host          = "https://localhost:9443"
  client_id     = "NPA"
  client_secret = "123456"
}

resource "wso2apim_application" "example" {
  name             = "foo-service"
  description      = "This is a foo service"
  throttling_policy = "Unlimited"
  token_type       = "JWT"
}

output "application_id" {
  value = wso2apim_application.example.id
}
```

## Installation

Add the following to your terraform configuration

```tf
terraform {
  required_providers {
    wso2apim = {
      source  = "floydspace/wso2apim"
      version = "~> 0.1.0"
    }
  }
}
```
