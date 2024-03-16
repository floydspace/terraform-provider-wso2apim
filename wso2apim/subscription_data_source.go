package wso2apim

import (
	"context"

	"github.com/floydspace/terraform-provider-wso2apim/apim"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &subscriptionDataSource{}
)

// NewSubscriptionDataSource is a helper function to simplify the provider implementation.
func NewSubscriptionDataSource() datasource.DataSource {
	return &subscriptionDataSource{}
}

// subscriptionDataSource is the data source implementation.
type subscriptionDataSource struct {
}

// subscriptionDataSourceModel maps the data source schema data.
type subscriptionDataSourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ApplicationID             types.String `tfsdk:"application_id"`
	ApiID                     types.String `tfsdk:"api_id"`
	ThrottlingPolicy          types.String `tfsdk:"throttling_policy"`
	RequestedThrottlingPolicy types.String `tfsdk:"requested_throttling_policy"`
	Status                    types.String `tfsdk:"status"`
}

// Metadata returns the data source type name.
func (d *subscriptionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subscription"
}

// Schema defines the schema for the data source.
func (d *subscriptionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a WSO2 API Manager Subscription",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Subscription ID.",
				Required:    true,
			},
			"application_id": schema.StringAttribute{
				Description: "Application ID.",
				Computed:    true,
			},
			"api_id": schema.StringAttribute{
				Description: "API ID.",
				Computed:    true,
			},
			"throttling_policy": schema.StringAttribute{
				Description: "Throttling policy.",
				Computed:    true,
			},
			"requested_throttling_policy": schema.StringAttribute{
				Description: "Requested throttling policy.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Subscription status.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *subscriptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state subscriptionDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	subscription, err := apim.GetSubscription(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading WSO2 API Manager Subscription",
			"Could not read subscription ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(subscription.SubscriptionID)
	state.ApplicationID = types.StringValue(subscription.ApplicationID)
	state.ApiID = types.StringValue(subscription.ApiID)
	state.ThrottlingPolicy = types.StringValue(subscription.ThrottlingPolicy)
	state.RequestedThrottlingPolicy = types.StringValue(subscription.RequestedThrottlingPolicy)
	state.Status = types.StringValue(subscription.Status)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
