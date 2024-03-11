package wso2apim

import (
	"context"
	"strings"
	"time"

	"github.com/floydspace/terraform-provider-wso2apim/apim"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &apiResource{}
	_ resource.ResourceWithImportState = &apiResource{}
	_ resource.ResourceWithConfigure   = &apiResource{}
)

// NewApiResource is a helper function to simplify the provider implementation.
func NewApiResource() resource.Resource {
	return &apiResource{}
}

// apiResource is the resource implementation.
type apiResource struct {
	config *wso2apimProviderModel
}

// apiResourceModel maps the resource schema data.
type apiResourceModel struct {
	ID              types.String                `tfsdk:"id"`
	Name            types.String                `tfsdk:"name"`
	Description     types.String                `tfsdk:"description"`
	Context         types.String                `tfsdk:"context"`
	Version         types.String                `tfsdk:"version"`
	Provider        types.String                `tfsdk:"api_provider"`
	Type            types.String                `tfsdk:"type"`
	LifeCycleStatus types.String                `tfsdk:"lifecycle_status"`
	HasThumbnail    types.Bool                  `tfsdk:"has_thumbnail"`
	Operations      []apiOperationResourceModel `tfsdk:"operations"`
	LastUpdated     types.String                `tfsdk:"last_updated"`
}

type apiOperationResourceModel struct {
	// ID     types.String `tfsdk:"id"`
	Target types.String `tfsdk:"target"`
	Verb   types.String `tfsdk:"verb"`
}

// Configure adds the provider configuration to the resource.
func (r *apiResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.config = req.ProviderData.(*wso2apimProviderModel)
}

// Metadata returns the resource type name.
func (r *apiResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api"
}

// Schema defines the schema for the resource.
func (r *apiResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a WSO2 API Manager Api.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Api ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the api.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the api.",
				Optional:    true,
			},
			"context": schema.StringAttribute{
				Description: "Context of the api.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Description: "Version of the api.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"api_provider": schema.StringAttribute{
				Description: "Provider of the api.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Type of the api.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("HTTP", "WS", "SOAPTOREST", "SOAP", "GRAPHQL", "WEBSUB", "SSE", "WEBHOOK", "ASYNC"),
				},
			},
			"lifecycle_status": schema.StringAttribute{
				Description: "LifeCycle status of the api.",
				Computed:    true,
			},
			"has_thumbnail": schema.BoolAttribute{
				Description: "Whether the api has a thumbnail.",
				Computed:    true,
			},
			"operations": schema.ListNestedAttribute{
				Description: "Operations of the api (Resources).",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// "id": schema.StringAttribute{
						// 	Description: "Operation ID.",
						// 	Computed:    true,
						// },
						"target": schema.StringAttribute{
							Description: "Operation target.",
							Optional:    true,
						},
						"verb": schema.StringAttribute{
							Description: "Operation verb.",
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"),
							},
						},
					},
				},
			},
			"last_updated": schema.StringAttribute{
				Description: "Last updated timestamp.",
				Computed:    true,
			},
		},
	}
}

// Create a new resource
func (r *apiResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan apiResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var description *string
	if !plan.Description.IsUnknown() {
		description = plan.Description.ValueStringPointer()
	}

	var apiProvider *string
	if !plan.Provider.IsUnknown() {
		apiProvider = plan.Provider.ValueStringPointer()
	}

	var apiType *string
	if !plan.Type.IsUnknown() {
		apiType = plan.Type.ValueStringPointer()
	}

	// Create operations
	var operations []apim.APIOperation
	for _, operation := range plan.Operations {
		operations = append(operations, apim.APIOperation{
			// ID:     operation.ID.ValueString(),
			Target: operation.Target.ValueString(),
			Verb:   operation.Verb.ValueString(),
		})
	}

	// Create new api
	api, err := apim.CreateAPI(&apim.APIReqBody{
		Name:        plan.Name.ValueString(),
		Description: description,
		Context:     plan.Context.ValueString(),
		Version:     plan.Version.ValueString(),
		Provider:    apiProvider,
		Type:        apiType,
		Operations:  operations,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating api",
			"Could not create api, unexpected error: "+err.Error(),
		)
		return
	}

	apiContext := api.Context
	if !r.config.ApiContextPrefix.IsUnknown() {
		apiContext = strings.Split(api.Context, r.config.ApiContextPrefix.ValueString())[1]
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(api.ID)
	plan.Name = types.StringValue(api.Name)
	plan.Description = types.StringValue(api.Description)
	plan.Context = types.StringValue(apiContext)
	plan.Version = types.StringValue(api.Version)
	plan.Provider = types.StringValue(api.Provider)
	plan.Type = types.StringValue(api.Type)
	plan.LifeCycleStatus = types.StringValue(api.LifeCycleStatus)
	plan.HasThumbnail = types.BoolValue(api.HasThumbnail)
	var planOperations []apiOperationResourceModel
	for _, operation := range api.Operations {
		planOperations = append(planOperations, apiOperationResourceModel{
			// ID:     types.StringValue(operation.ID),
			Target: types.StringValue(operation.Target),
			Verb:   types.StringValue(operation.Verb),
		})
	}
	plan.Operations = planOperations
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *apiResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state apiResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed api value from WSO2 API Manager
	api, err := apim.GetAPI(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading WSO2 API Manager Api",
			"Could not read api ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	apiContext := api.Context
	if !r.config.ApiContextPrefix.IsUnknown() {
		apiContext = strings.Split(api.Context, r.config.ApiContextPrefix.ValueString())[1]
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(api.Name)
	state.Description = types.StringValue(api.Description)
	state.Context = types.StringValue(apiContext)
	state.Version = types.StringValue(api.Version)
	state.Provider = types.StringValue(api.Provider)
	state.Type = types.StringValue(api.Type)
	state.LifeCycleStatus = types.StringValue(api.LifeCycleStatus)
	state.HasThumbnail = types.BoolValue(api.HasThumbnail)
	var operations []apiOperationResourceModel
	for _, operation := range api.Operations {
		operations = append(operations, apiOperationResourceModel{
			// ID:     types.StringValue(operation.ID),
			Target: types.StringValue(operation.Target),
			Verb:   types.StringValue(operation.Verb),
		})
	}
	state.Operations = operations

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *apiResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan apiResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var description *string
	if !plan.Description.IsUnknown() {
		description = plan.Description.ValueStringPointer()
	}

	var apiType *string
	if !plan.Type.IsUnknown() {
		apiType = plan.Type.ValueStringPointer()
	}

	var operations []apim.APIOperation
	for _, operation := range plan.Operations {
		operations = append(operations, apim.APIOperation{
			// ID:     operation.ID.ValueString(),
			Target: operation.Target.ValueString(),
			Verb:   operation.Verb.ValueString(),
		})
	}

	// Create new api
	api, err := apim.UpdateAPI(plan.ID.ValueString(), &apim.APIReqBody{
		Description: description,
		Type:        apiType,
		Operations:  operations,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating api",
			"Could not create api, unexpected error: "+err.Error(),
		)
		return
	}

	apiContext := api.Context
	if !r.config.ApiContextPrefix.IsUnknown() {
		apiContext = strings.Split(api.Context, r.config.ApiContextPrefix.ValueString())[1]
	}

	// Update resource state with updated items and timestamp
	plan.ID = types.StringValue(api.ID)
	plan.Name = types.StringValue(api.Name)
	plan.Description = types.StringValue(api.Description)
	plan.Context = types.StringValue(apiContext)
	plan.Version = types.StringValue(api.Version)
	plan.Provider = types.StringValue(api.Provider)
	plan.Type = types.StringValue(api.Type)
	plan.LifeCycleStatus = types.StringValue(api.LifeCycleStatus)
	plan.HasThumbnail = types.BoolValue(api.HasThumbnail)
	var planOperations []apiOperationResourceModel
	for _, operation := range api.Operations {
		planOperations = append(planOperations, apiOperationResourceModel{
			// ID:     types.StringValue(operation.ID),
			Target: types.StringValue(operation.Target),
			Verb:   types.StringValue(operation.Verb),
		})
	}
	plan.Operations = planOperations
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *apiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state apiResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing api
	err := apim.DeleteAPI(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting WSO2 API Manager Api",
			"Could not delete api, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *apiResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
