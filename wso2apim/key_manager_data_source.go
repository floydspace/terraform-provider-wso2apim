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
	_ datasource.DataSource = &keyManagerDataSource{}
)

// NewKeyManagerDataSource is a helper function to simplify the provider implementation.
func NewKeyManagerDataSource() datasource.DataSource {
	return &keyManagerDataSource{}
}

// keyManagerDataSource is the data source implementation.
type keyManagerDataSource struct {
}

// keyManagerDataSourceModel maps the data source schema data.
type keyManagerDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Type                types.String `tfsdk:"type"`
	DisplayName         types.String `tfsdk:"display_name"`
	Description         types.String `tfsdk:"description"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	AvailableGrantTypes []string     `tfsdk:"available_grant_types"`
	TokenEndpoint       types.String `tfsdk:"token_endpoint"`
	RevokeEndpoint      types.String `tfsdk:"revoke_endpoint"`
}

// Metadata returns the data source type name.
func (d *keyManagerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_manager"
}

// Schema defines the schema for the data source.
func (d *keyManagerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a WSO2 API Manager KeyManager",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Key Manager ID.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Key Manager Name.",
				Optional:    true,
			},
			"type": schema.StringAttribute{
				Description: "Key Manager Type.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "Key Manager Display Name.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Key Manager Description.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Key Manager Enabled.",
				Computed:    true,
			},
			"available_grant_types": schema.ListAttribute{
				Description: "Key Manager Available Grant Types.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"token_endpoint": schema.StringAttribute{
				Description: "Key Manager Token Endpoint.",
				Computed:    true,
			},
			"revoke_endpoint": schema.StringAttribute{
				Description: "Key Manager Revoke Endpoint.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *keyManagerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state keyManagerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var keyManager *apim.KeyManagerSearchInfo

	if !state.ID.IsUnknown() && !state.ID.IsNull() {
		km, err := apim.GetKeyManager(state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading WSO2 API Manager KeyManager",
				"Could not read keyManager ID "+state.ID.ValueString()+": "+err.Error(),
			)
			return
		}
		keyManager = km
	} else if !state.Name.IsUnknown() && !state.Name.IsNull() {
		km, err := apim.SearchKeyManager(state.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading WSO2 API Manager KeyManager",
				"Could not read keyManager Name "+state.Name.ValueString()+": "+err.Error(),
			)
			return
		}
		keyManager = km
	} else {
		resp.Diagnostics.AddError(
			"Error Reading WSO2 API Manager KeyManager",
			"Either id or name must be provided",
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(keyManager.ID)
	state.Name = types.StringValue(keyManager.Name)
	state.Type = types.StringValue(keyManager.Type)
	state.DisplayName = types.StringValue(keyManager.DisplayName)
	state.Description = types.StringValue(keyManager.Description)
	state.Enabled = types.BoolValue(keyManager.Enabled)
	state.AvailableGrantTypes = keyManager.AvailableGrantTypes
	state.TokenEndpoint = types.StringValue(keyManager.TokenEndpoint)
	state.RevokeEndpoint = types.StringValue(keyManager.RevokeEndpoint)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
