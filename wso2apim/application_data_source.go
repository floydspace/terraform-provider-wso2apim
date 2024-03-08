package wso2apim

import (
	"context"

	"github.com/floydspace/terraform-provider-wso2apim/apim"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &applicationDataSource{}
)

// NewApplicationDataSource is a helper function to simplify the provider implementation.
func NewApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

// applicationDataSource is the data source implementation.
type applicationDataSource struct {
}

// applicationDataSourceModel maps the data source schema data.
type applicationDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ThrottlingPolicy  types.String `tfsdk:"throttling_policy"`
	Description       types.String `tfsdk:"description"`
	Status            types.String `tfsdk:"status"`
	SubscriptionCount types.Int64  `tfsdk:"subscription_count"`
	Attributes        types.Map    `tfsdk:"attributes"`
	Owner             types.String `tfsdk:"owner"`
	TokenType         types.String `tfsdk:"token_type"`
}

// Metadata returns the data source type name.
func (d *applicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

// Schema defines the schema for the data source.
func (d *applicationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a WSO2 API Manager Application",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Application ID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the application.",
				Computed:    true,
			},
			"throttling_policy": schema.StringAttribute{
				Description: "Application throttling policy.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the application.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the application.",
				Computed:    true,
			},
			"subscription_count": schema.Int64Attribute{
				Description: "Number of subscriptions to the application.",
				Computed:    true,
			},
			"attributes": schema.MapAttribute{
				Description: "Attributes of the application.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"owner": schema.StringAttribute{
				Description: "Owner of the application.",
				Computed:    true,
			},
			"token_type": schema.StringAttribute{
				Description: "Token type of the application.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state applicationDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	application, err := apim.GetApplication(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading WSO2 API Manager Application",
			"Could not read application ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(application.ApplicationID)
	state.Name = types.StringValue(application.Name)
	state.ThrottlingPolicy = types.StringValue(application.ThrottlingPolicy)
	state.Description = types.StringValue(application.Description)
	state.Status = types.StringValue(application.Status)
	state.SubscriptionCount = types.Int64Value(int64(application.SubscriptionCount))
	appAttributes := make(map[string]attr.Value)
	for k, v := range application.Attributes {
		appAttributes[k] = types.StringValue(v)
	}
	state.Attributes = types.MapValueMust(types.StringType, appAttributes)
	state.Owner = types.StringValue(application.Owner)
	state.TokenType = types.StringValue(application.TokenType)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
