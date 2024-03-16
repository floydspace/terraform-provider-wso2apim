package wso2apim

import (
	"context"
	"strings"
	"time"

	"github.com/floydspace/terraform-provider-wso2apim/apim"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &applicationKeyMappingResource{}
	_ resource.ResourceWithImportState = &applicationKeyMappingResource{}
)

// NewApplicationKeyMappingResource is a helper function to simplify the provider implementation.
func NewApplicationKeyMappingResource() resource.Resource {
	return &applicationKeyMappingResource{}
}

// applicationKeyMappingResource is the resource implementation.
type applicationKeyMappingResource struct {
}

// applicationKeyMappingResourceModel maps the resource schema data.
type applicationKeyMappingResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	ApplicationID       types.String `tfsdk:"application_id"`
	KeyType             types.String `tfsdk:"key_type"`
	KeyManager          types.String `tfsdk:"key_manager"`
	ConsumerKey         types.String `tfsdk:"consumer_key"`
	ConsumerSecret      types.String `tfsdk:"consumer_secret"`
	SupportedGrantTypes []string     `tfsdk:"supported_grant_types"`
	// CallbackURL         types.String `tfsdk:"callback_url"`
	ValidityTime types.Int64  `tfsdk:"validity_time"`
	Scopes       []string     `tfsdk:"scopes"`
	KeyState     types.String `tfsdk:"key_state"`
	LastUpdated  types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *applicationKeyMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_key_mapping"
}

// Schema defines the schema for the resource.
func (r *applicationKeyMappingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a WSO2 API Manager Application Keys.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Application Key Mapping ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				Description: "Application ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_type": schema.StringAttribute{
				Description: "Application Key Type.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("PRODUCTION", "SANDBOX"),
				},
			},
			"key_manager": schema.StringAttribute{
				Description: "Key Manager.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"consumer_key": schema.StringAttribute{
				Description: "Consumer Key.",
				Computed:    true,
			},
			"consumer_secret": schema.StringAttribute{
				Description: "Consumer Secret.",
				Computed:    true,
				Sensitive:   true,
			},
			"supported_grant_types": schema.ListAttribute{
				Description: "Supported Grant Types.",
				ElementType: types.StringType,
				Required:    true,
			},
			// "callback_url": schema.StringAttribute{
			// 	Description: "Callback URL.",
			// 	Optional:    true,
			// },
			"validity_time": schema.Int64Attribute{
				Description: "Validity Time.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
				Default: int64default.StaticInt64(3600),
			},
			"scopes": schema.ListAttribute{
				Description: "Scopes.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
					listplanmodifier.RequiresReplace(),
				},
				Default: listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{
					types.StringValue("default"),
				})),
			},
			"key_state": schema.StringAttribute{
				Description: "Application Key State.",
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
func (r *applicationKeyMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan applicationKeyMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new applicationKeyMapping
	applicationKeyMapping, err := apim.GenerateKeys(plan.ApplicationID.ValueString(), &apim.ApplicationKeyGenerateRequest{
		KeyType:                 plan.KeyType.ValueString(),
		KeyManager:              plan.KeyManager.ValueString(),
		GrantTypesToBeSupported: plan.SupportedGrantTypes,
		// CallbackURL:             plan.CallbackURL.ValueString(),
		ValidityTime: plan.ValidityTime.ValueInt64(),
		Scopes:       plan.Scopes,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating application keys",
			"Could not create applicationKeyMapping, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(applicationKeyMapping.ID)
	plan.KeyType = types.StringValue(applicationKeyMapping.KeyType)
	plan.KeyManager = types.StringValue(applicationKeyMapping.KeyManager)
	plan.ConsumerKey = types.StringValue(applicationKeyMapping.ConsumerKey)
	plan.ConsumerSecret = types.StringValue(applicationKeyMapping.ConsumerSecret)
	plan.SupportedGrantTypes = applicationKeyMapping.SupportedGrantTypes
	// if applicationKeyMapping.CallbackURL != nil {
	// 	plan.CallbackURL = types.StringPointerValue(applicationKeyMapping.CallbackURL)
	// }
	plan.KeyState = types.StringValue(applicationKeyMapping.KeyState)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *applicationKeyMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state applicationKeyMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed applicationKeyMapping value from WSO2 API Manager
	applicationKeyMapping, err := apim.GetApplicationKeys(state.ApplicationID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading WSO2 API Manager Application Keys",
			"Could not read applicationKeyMapping ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.KeyType = types.StringValue(applicationKeyMapping.KeyType)
	state.KeyManager = types.StringValue(applicationKeyMapping.KeyManager)
	state.ConsumerKey = types.StringValue(applicationKeyMapping.ConsumerKey)
	state.ConsumerSecret = types.StringValue(applicationKeyMapping.ConsumerSecret)
	state.SupportedGrantTypes = applicationKeyMapping.SupportedGrantTypes
	// if applicationKeyMapping.CallbackURL != nil {
	// 	state.CallbackURL = types.StringPointerValue(applicationKeyMapping.CallbackURL)
	// }
	state.KeyState = types.StringValue(applicationKeyMapping.KeyState)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *applicationKeyMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan applicationKeyMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mapping, err := apim.UpdateApplicationKeys(plan.ApplicationID.ValueString(), plan.ID.ValueString(), &apim.ApplicationKeyGenerateRequest{
		GrantTypesToBeSupported: plan.SupportedGrantTypes,
		// CallbackURL:             plan.CallbackURL.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating application keys",
			"Could not update application key mapping, unexpected error: "+err.Error(),
		)
		return
	}

	// keys, err := apim.RegenerateKeys(plan.ApplicationID.ValueString(), plan.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error regenerating application keys",
	// 		"Could not create application key mapping, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	// Update resource state with updated items and timestamp
	// plan.ConsumerKey = types.StringValue(keys.ConsumerKey)
	// plan.ConsumerSecret = types.StringValue(keys.ConsumerSecret)
	plan.SupportedGrantTypes = mapping.SupportedGrantTypes
	// if mapping.CallbackURL != nil {
	// 	plan.CallbackURL = types.StringPointerValue(mapping.CallbackURL)
	// }
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *applicationKeyMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state applicationKeyMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing applicationKeyMapping
	err := apim.CleanupKeys(state.ApplicationID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Cleaning up WSO2 API Manager Application Keys",
			"Could not delete application key mapping, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *applicationKeyMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	parts := strings.Split(req.ID, ",")

	if len(parts) < 2 {
		resp.Diagnostics.AddError(
			"Error importing item",
			"Could not import item, unexpected error (ID should be in the format <application_id>,<id>): "+req.ID,
		)
		return
	}

	appID := parts[0]
	ID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), appID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ID)...)
}
