package wso2apim

import (
	"context"
	"fmt"
	"os"

	"github.com/floydspace/terraform-provider-wso2apim/apim"
	"github.com/floydspace/terraform-provider-wso2apim/token"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/wso2/openservicebroker-apim/pkg/client"
	apimCfg "github.com/wso2/openservicebroker-apim/pkg/config"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &wso2apimProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &wso2apimProvider{
			version: version,
		}
	}
}

// wso2apimProvider is the provider implementation.
type wso2apimProvider struct {
	version string
}

// wso2apimProviderModel maps provider schema data to a Go type.
type wso2apimProviderModel struct {
	Host             types.String `tfsdk:"host"`
	Username         types.String `tfsdk:"username"`
	Password         types.String `tfsdk:"password"`
	ApiContextPrefix types.String `tfsdk:"api_context_prefix"`
}

// Metadata returns the provider type name.
func (p *wso2apimProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "wso2apim"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *wso2apimProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with WSO2 API Manager.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "WSO2 API Manager Hostname. May also be provided via the WSO2_APIM_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "WSO2 API Manager Username. May also be provided via the WSO2_APIM_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "WSO2 API Manager Password. May also be provided via the WSO2_APIM_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"api_context_prefix": schema.StringAttribute{
				Description: "WSO2 API Manager API Context Prefix.",
				Optional:    true,
			},
		},
	}
}

// Configure prepares a wso2apim API client for data sources and resources.
func (p *wso2apimProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring WSO2 API Manager client")

	// Retrieve provider data from configuration
	var config wso2apimProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown WSO2 API Manager Host",
			"The provider cannot create the WSO2 API Manager Consumer client as there is an unknown configuration value for the WSO2 API Manager Host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WSO2_APIM_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown WSO2 API Manager Username",
			"The provider cannot create the WSO2 API Manager Consumer client as there is an unknown configuration value for the WSO2 API Manager Username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WSO2_APIM_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown WSO2 API Manager Password",
			"The provider cannot create the WSO2 API Manager Consumer client as there is an unknown configuration value for the WSO2 API Manager Password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WSO2_APIM_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("WSO2_APIM_HOST")
	username := os.Getenv("WSO2_APIM_USERNAME")
	password := os.Getenv("WSO2_APIM_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing WSO2 API Manager Host",
			"The provider cannot create the WSO2 API Manager client as there is a missing or empty value for the WSO2 Host. "+
				"Set the host value in the configuration or use the WSO2_APIM_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing WSO2 API Manager Username",
			"The provider cannot create the WSO2 API Manager client as there is a missing or empty value for the WSO2 API Manager Username. "+
				"Set the username value in the configuration or use the WSO2_APIM_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing WSO2 API Manager Password",
			"The provider cannot create the WSO2 API Manager client as there is a missing or empty value for the WSO2 API Manager Password. "+
				"Set the password value in the configuration or use the WSO2_APIM_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "wso2apim_host", host)
	ctx = tflog.SetField(ctx, "wso2apim_username", username)
	ctx = tflog.SetField(ctx, "wso2apim_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "wso2apim_password")

	tflog.Debug(ctx, "Creating WSO2 API Manager client")

	apimConf := apim.APIM{
		Username:                         username,
		Password:                         password,
		TokenEndpoint:                    host + "/oauth2",
		DynamicClientEndpoint:            host,
		DynamicClientRegistrationContext: "/client-registration/v0.17/register",
		PublisherEndpoint:                host,
		PublisherAPIContext:              "/api/am/publisher/v1/apis",
		StoreEndpoint:                    host,
		StoreApplicationContext:          "/api/am/store/v1/applications",
		StoreKeyManagerContext:           "/api/am/store/v1/key-managers",
		StoreSubscriptionContext:         "/api/am/store/v1/subscriptions",
		StoreMultipleSubscriptionContext: "/api/am/store/v1/subscriptions/multiple",
	}

	client.Configure(&apimCfg.Client{
		InsecureCon: true,
		MinBackOff:  1,
		MaxBackOff:  60,
		Timeout:     30,
		MaxRetries:  3,
	})

	// Initialize Token manager.
	tManager := &token.PasswordRefreshTokenGrantManager{
		TokenEndpoint:                    apimConf.TokenEndpoint,
		DynamicClientEndpoint:            apimConf.DynamicClientEndpoint,
		DynamicClientRegistrationContext: apimConf.DynamicClientRegistrationContext,
		UserName:                         apimConf.Username,
		Password:                         apimConf.Password,
	}
	tManager.Init([]string{
		token.ScopeSubscribe,
		token.ScopeAPIView,
		token.ScopeAPICreate,
		token.ScopeAppPublish,
		token.ScopeAPIDelete,
		token.ScopeAppManage,
	})

	defer func() {
		if r := recover(); r != nil {
			resp.Diagnostics.AddError(
				"Unable to Create WSO2 API Manager Client",
				"An unexpected error occurred when creating the WSO2 API Manager client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"Client Error: "+fmt.Sprintf("%v", r),
			)
		}
	}()

	// Create a new WSO2 client using the configuration values
	apim.Init(tManager, apimConf)

	resp.ResourceData = &config

	tflog.Info(ctx, "Configured WSO2 API Manager client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *wso2apimProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewApiDataSource,
		NewKeyManagerDataSource,
		NewApplicationDataSource,
		NewSubscriptionDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *wso2apimProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApiResource,
		NewApplicationResource,
		NewSubscriptionResource,
	}
}
