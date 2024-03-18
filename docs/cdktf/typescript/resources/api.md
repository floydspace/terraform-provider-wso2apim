---

<!-- Please do not edit this file, it is generated. -->
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "wso2apim_api Resource - wso2apim"
subcategory: ""
description: |-
  Manages a WSO2 API Manager Api.
---

# wso2apim_api (Resource)

Manages a WSO2 API Manager Api.

## Example Usage

```typescript
// DO NOT EDIT. Code generated by 'cdktf convert' - Please report bugs at https://cdk.tf/bug
import { Construct } from "constructs";
import { TerraformStack } from "cdktf";
/*
 * Provider bindings are generated by running `cdktf get`.
 * See https://cdk.tf/provider-generation for more details.
 */
import { Api } from "./.gen/providers/wso2apim/api";
class MyConvertedCode extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);
    new Api(this, "example", {
      context: "/foo",
      description: "This is a foo API",
      name: "foo-api",
      version: "v1",
    });
  }
}

```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `context` (String) Context of the api.
- `name` (String) Name of the api.
- `version` (String) Version of the api.

### Optional

- `apiProvider` (String) Provider of the api.
- `description` (String) Description of the api.
- `endpointConfig` (Attributes) Endpoint configuration of the api. (see [below for nested schema](#nestedatt--endpoint_config))
- `operations` (Attributes List) Operations of the api (Resources). (see [below for nested schema](#nestedatt--operations))
- `policies` (List of String) Policies of the api.
- `type` (String) Type of the api.

### Read-Only

- `hasThumbnail` (Boolean) Whether the api has a thumbnail.
- `id` (String) Api ID.
- `lastUpdated` (String) Last updated timestamp.
- `lifecycleStatus` (String) LifeCycle status of the api.

<a id="nestedatt--endpoint_config"></a>
### Nested Schema for `endpointConfig`

Optional:

- `endpointType` (String) Endpoint type.
- `production_endpoints` (Attributes) Sandbox endpoints. (see [below for nested schema](#nestedatt--endpoint_config--production_endpoints))
- `sandbox_endpoints` (Attributes) Sandbox endpoints. (see [below for nested schema](#nestedatt--endpoint_config--sandbox_endpoints))

<a id="nestedatt--endpoint_config--production_endpoints"></a>
### Nested Schema for `endpoint_config.production_endpoints`

Optional:

- `url` (String) Sandbox endpoint URL.


<a id="nestedatt--endpoint_config--sandbox_endpoints"></a>
### Nested Schema for `endpoint_config.sandbox_endpoints`

Optional:

- `url` (String) Sandbox endpoint URL.



<a id="nestedatt--operations"></a>
### Nested Schema for `operations`

Optional:

- `target` (String) Operation target.
- `verb` (String) Operation verb.

## Import

Import is supported using the following syntax:

```shell
# Api can be imported by specifying the id of the api
terraform import wso2apim_api.example 00000000-0000-0000-0000-000000000000
```

<!-- cache-key: cdktf-0.20.4 input-539b678364f745518e337adb74e514e6acb0c3126b5f70d2d807740b9eccd788 -->