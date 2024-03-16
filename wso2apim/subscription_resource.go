package wso2apim

import (
	"context"
	"time"

	"github.com/floydspace/terraform-provider-wso2apim/apim"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &subscriptionResource{}
	_ resource.ResourceWithImportState = &subscriptionResource{}
)

// NewSubscriptionResource is a helper function to simplify the provider implementation.
func NewSubscriptionResource() resource.Resource {
	return &subscriptionResource{}
}

// subscriptionResource is the resource implementation.
type subscriptionResource struct {
}

// subscriptionResourceModel maps the resource schema data.
type subscriptionResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	ApplicationID             types.String `tfsdk:"application_id"`
	ApiID                     types.String `tfsdk:"api_id"`
	ThrottlingPolicy          types.String `tfsdk:"throttling_policy"`
	RequestedThrottlingPolicy types.String `tfsdk:"requested_throttling_policy"`
	Status                    types.String `tfsdk:"status"`
	LastUpdated               types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *subscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subscription"
}

// Schema defines the schema for the resource.
func (r *subscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a WSO2 API Manager Subscription.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Subscription ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				Description: "Application ID.",
				Required:    true,
			},
			"api_id": schema.StringAttribute{
				Description: "API ID.",
				Optional:    true,
			},
			"throttling_policy": schema.StringAttribute{
				Description: "Throttling policy.",
				Required:    true,
			},
			"requested_throttling_policy": schema.StringAttribute{
				Description: "Requested throttling policy.",
				Optional:    true,
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Subscription status.",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Description: "Last updated timestamp.",
				Computed:    true,
			},
		},
	}
}

// Create a new resource
func (r *subscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan subscriptionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new subscription
	subscription, err := apim.CreateSubscription(&apim.SubscriptionReq{
		ApplicationID:             plan.ApplicationID.ValueString(),
		ApiID:                     plan.ApiID.ValueString(),
		ThrottlingPolicy:          plan.ThrottlingPolicy.ValueString(),
		RequestedThrottlingPolicy: plan.RequestedThrottlingPolicy.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating subscription",
			"Could not create subscription, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(subscription.SubscriptionID)
	plan.ApplicationID = types.StringValue(subscription.ApplicationID)
	plan.ApiID = types.StringValue(subscription.ApiID)
	plan.ThrottlingPolicy = types.StringValue(subscription.ThrottlingPolicy)
	plan.RequestedThrottlingPolicy = types.StringValue(subscription.RequestedThrottlingPolicy)
	plan.Status = types.StringValue(subscription.Status)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *subscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state subscriptionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed subscription value from WSO2 API Manager
	subscription, err := apim.GetSubscription(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading WSO2 API Manager Subscription",
			"Could not read subscription ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ApplicationID = types.StringValue(subscription.ApplicationID)
	state.ApiID = types.StringValue(subscription.ApiID)
	state.ThrottlingPolicy = types.StringValue(subscription.ThrottlingPolicy)
	state.RequestedThrottlingPolicy = types.StringValue(subscription.RequestedThrottlingPolicy)
	state.Status = types.StringValue(subscription.Status)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *subscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan subscriptionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new subscription
	subscription, err := apim.UpdateSubscription(plan.ID.ValueString(), &apim.SubscriptionReq{
		ApplicationID:             plan.ApplicationID.ValueString(),
		ApiID:                     plan.ApiID.ValueString(),
		ThrottlingPolicy:          plan.ThrottlingPolicy.ValueString(),
		RequestedThrottlingPolicy: plan.RequestedThrottlingPolicy.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating subscription",
			"Could not create subscription, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.ID = types.StringValue(subscription.SubscriptionID)
	plan.ApplicationID = types.StringValue(subscription.ApplicationID)
	plan.ApiID = types.StringValue(subscription.ApiID)
	plan.ThrottlingPolicy = types.StringValue(subscription.ThrottlingPolicy)
	plan.RequestedThrottlingPolicy = types.StringValue(subscription.RequestedThrottlingPolicy)
	plan.Status = types.StringValue(subscription.Status)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *subscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state subscriptionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing subscription
	err := apim.UnSubscribe(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting WSO2 API Manager Subscription",
			"Could not delete subscription, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *subscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
