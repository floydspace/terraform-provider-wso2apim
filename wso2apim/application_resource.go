package wso2apim

import (
	"context"
	"time"

	"github.com/floydspace/terraform-provider-wso2apim/apim"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &applicationResource{}
	_ resource.ResourceWithImportState = &applicationResource{}
)

// NewApplicationResource is a helper function to simplify the provider implementation.
func NewApplicationResource() resource.Resource {
	return &applicationResource{}
}

// applicationResource is the resource implementation.
type applicationResource struct {
}

// applicationResourceModel maps the resource schema data.
type applicationResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ThrottlingPolicy  types.String `tfsdk:"throttling_policy"`
	Description       types.String `tfsdk:"description"`
	Status            types.String `tfsdk:"status"`
	SubscriptionCount types.Int64  `tfsdk:"subscription_count"`
	Attributes        types.Map    `tfsdk:"attributes"`
	Owner             types.String `tfsdk:"owner"`
	TokenType         types.String `tfsdk:"token_type"`
	LastUpdated       types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

// Schema defines the schema for the resource.
func (r *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a WSO2 API Manager Application.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Application ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the application.",
				Required:    true,
			},
			"throttling_policy": schema.StringAttribute{
				Description: "Application throttling policy.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the application.",
				Optional:    true,
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
				Optional:    true,
			},
			"owner": schema.StringAttribute{
				Description: "Owner of the application.",
				Computed:    true,
			},
			"token_type": schema.StringAttribute{
				Description: "Token type of the application.",
				Required:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Last updated timestamp.",
				Computed:    true,
			},
		},
	}
}

// Create a new resource
func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan applicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var attributes map[string]string
	diags = plan.Attributes.ElementsAs(ctx, &attributes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new application
	application, err := apim.CreateApplication(&apim.ApplicationCreateReq{
		Name:             plan.Name.ValueString(),
		TokenType:        plan.TokenType.ValueString(),
		ThrottlingPolicy: plan.ThrottlingPolicy.ValueString(),
		Description:      plan.Description.ValueString(),
		Attributes:       attributes,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating application",
			"Could not create application, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(application.ApplicationID)
	plan.Name = types.StringValue(application.Name)
	plan.ThrottlingPolicy = types.StringValue(application.ThrottlingPolicy)
	plan.Description = types.StringValue(application.Description)
	plan.Status = types.StringValue(application.Status)
	plan.SubscriptionCount = types.Int64Value(int64(application.SubscriptionCount))
	appAttributes := make(map[string]attr.Value)
	for k, v := range application.Attributes {
		appAttributes[k] = types.StringValue(v)
	}
	plan.Attributes = types.MapValueMust(types.StringType, appAttributes)
	plan.Owner = types.StringValue(application.Owner)
	plan.TokenType = types.StringValue(application.TokenType)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state applicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed application value from WSO2 API Manager
	application, err := apim.GetApplication(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading WSO2 API Manager Application",
			"Could not read application ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
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

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan applicationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing application
	err := apim.DeleteApplication(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting WSO2 API Manager Application",
			"Could not delete application, unexpected error: "+err.Error(),
		)
		return
	}

	var attributes map[string]string
	diags = plan.Attributes.ElementsAs(ctx, &attributes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new application
	application, err := apim.CreateApplication(&apim.ApplicationCreateReq{
		Name:             plan.Name.ValueString(),
		TokenType:        plan.TokenType.ValueString(),
		ThrottlingPolicy: plan.ThrottlingPolicy.ValueString(),
		Description:      plan.Description.ValueString(),
		Attributes:       attributes,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating application",
			"Could not create application, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.ID = types.StringValue(application.ApplicationID)
	plan.Name = types.StringValue(application.Name)
	plan.ThrottlingPolicy = types.StringValue(application.ThrottlingPolicy)
	plan.Description = types.StringValue(application.Description)
	plan.Status = types.StringValue(application.Status)
	plan.SubscriptionCount = types.Int64Value(int64(application.SubscriptionCount))
	appAttributes := make(map[string]attr.Value)
	for k, v := range application.Attributes {
		appAttributes[k] = types.StringValue(v)
	}
	plan.Attributes = types.MapValueMust(types.StringType, appAttributes)
	plan.Owner = types.StringValue(application.Owner)
	plan.TokenType = types.StringValue(application.TokenType)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state applicationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing application
	err := apim.DeleteApplication(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting WSO2 API Manager Application",
			"Could not delete application, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *applicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
