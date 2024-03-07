package main

import (
	"context"

	"github.com/floydspace/terraform-provider-wso2apim/wso2apim"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name wso2apim

func main() {
	providerserver.Serve(context.Background(), wso2apim.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/floydspace/wso2apim",
	})
}
