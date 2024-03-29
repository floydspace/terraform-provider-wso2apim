---

<!-- Please do not edit this file, it is generated. -->
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "wso2apim_subscription Data Source - wso2apim"
subcategory: ""
description: |-
  Fetches a WSO2 API Manager Subscription
---

# wso2apim_subscription (Data Source)

Fetches a WSO2 API Manager Subscription

## Example Usage

```python
# DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
from constructs import Construct
from cdktf import TerraformStack
#
# Provider bindings are generated by running `cdktf get`.
# See https://cdk.tf/provider-generation for more details.
#
from imports.wso2apim.data_wso2_apim_subscription import DataWso2ApimSubscription
class MyConvertedCode(TerraformStack):
    def __init__(self, scope, name):
        super().__init__(scope, name)
        DataWso2ApimSubscription(self, "example",
            id="00000000-0000-0000-0000-000000000000"
        )
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Subscription ID.

### Read-Only

- `api_id` (String) API ID.
- `application_id` (String) Application ID.
- `requested_throttling_policy` (String) Requested throttling policy.
- `status` (String) Subscription status.
- `throttling_policy` (String) Throttling policy.

<!-- cache-key: cdktf-0.20.4 input-68e907a36c8bb0c0708f1e96cd82fd45bc0867003b443aaf2c5f28cfdba0f072 -->