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
	_ datasource.DataSource = &apiDataSource{}
)

// NewApiDataSource is a helper function to simplify the provider implementation.
func NewApiDataSource() datasource.DataSource {
	return &apiDataSource{}
}

// apiDataSource is the data source implementation.
type apiDataSource struct {
}

// apiDataSourceModel maps the data source schema data.
type apiDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Context         types.String `tfsdk:"context"`
	Version         types.String `tfsdk:"version"`
	Provider        types.String `tfsdk:"api_provider"`
	Type            types.String `tfsdk:"type"`
	LifeCycleStatus types.String `tfsdk:"lifecycle_status"`
	HasThumbnail    types.Bool   `tfsdk:"has_thumbnail"`
}

// Metadata returns the data source type name.
func (d *apiDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api"
}

// Schema defines the schema for the data source.
func (d *apiDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a WSO2 API Manager Api",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Api ID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the api.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the api.",
				Computed:    true,
			},
			"context": schema.StringAttribute{
				Description: "Context of the api.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "Version of the api.",
				Computed:    true,
			},
			"api_provider": schema.StringAttribute{
				Description: "Provider of the api.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the api.",
				Computed:    true,
			},
			"lifecycle_status": schema.StringAttribute{
				Description: "LifeCycle status of the api.",
				Computed:    true,
			},
			"has_thumbnail": schema.BoolAttribute{
				Description: "Whether the api has a thumbnail.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *apiDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state apiDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	api, err := apim.GetAPI(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading WSO2 API Manager Api",
			"Could not read api ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(api.ID)
	state.Name = types.StringValue(api.Name)
	state.Description = types.StringValue(api.Description)
	state.Context = types.StringValue(api.Context)
	state.Version = types.StringValue(api.Version)
	state.Provider = types.StringValue(api.Provider)
	state.Type = types.StringValue(api.Type)
	state.LifeCycleStatus = types.StringValue(api.LifeCycleStatus)
	state.HasThumbnail = types.BoolValue(api.HasThumbnail)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
