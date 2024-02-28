package authorization

import (
	"context"
	"fmt"
	"slices"

	"github.com/kiwicom/terraform-provider-montecarlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ServiceAccountResource{}

// This resource cannot be imported, since Token cannot be retrieved from Monte Carlo API.
// var _ resource.ResourceWithImportState = &ServiceAccountResource{}

// To simplify provider implementations, a named function can be created with the resource implementation.
func NewServiceAccountResource() resource.Resource {
	return &ServiceAccountResource{}
}

// ServiceAccountResource defines the resource implementation.
type ServiceAccountResource struct {
	client client.MonteCarloClient
}

// ServiceAccountResourceModel describes the resource data model according to its Schema.
type ServiceAccountResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Token       types.String `tfsdk:"token"`
	Description types.String `tfsdk:"description"`
}

func (r *ServiceAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account"
}

func (r *ServiceAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token": schema.StringAttribute{
				Computed:  true,
				Optional:  false,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
			},
		},
	}
}

func (r *ServiceAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := common.Configure(req)
	resp.Diagnostics.Append(diags...)
	r.client = client
}

func (r *ServiceAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServiceAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResult := client.CreateOrUpdateServiceApiToken{}
	variables := map[string]interface{}{
		"tokenId":          (*string)(nil),
		"comment":          data.Description.ValueString(),
		"displayName":      (*string)(nil),
		"expirationInDays": (*int)(nil),
		"groups":           (*[]string)(nil),
	}

	if err := r.client.Mutate(ctx, &createResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'CreateOrUpdateServiceApiToken' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	data.Id = types.StringValue(createResult.CreateOrUpdateServiceApiToken.AccessToken.Id)
	data.Token = types.StringValue(createResult.CreateOrUpdateServiceApiToken.AccessToken.Token)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServiceAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var index int
	type AccessKeyIndexEnum string
	readResult := client.GetTokenMetadata{}
	variables := map[string]interface{}{
		"index":             (AccessKeyIndexEnum)("account"),
		"isServiceApiToken": true,
	}

	if err := r.client.Query(ctx, &readResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'GetTokenMetadata' query result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	if index = slices.IndexFunc(readResult.GetTokenMetadata, func(token client.TokenMetadata) bool {
		return token.Id == data.Id.ValueString()
	}); index < 0 {
		to_print := fmt.Sprintf("Token [ID: %s] not found", data.Id.ValueString())
		resp.Diagnostics.AddWarning(to_print, "")
		resp.State.RemoveResource(ctx)
		return
	}

	data.Description = types.StringValue(readResult.GetTokenMetadata[index].Comment)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServiceAccountResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateResult := client.CreateOrUpdateServiceApiToken{}
	variables := map[string]interface{}{
		"tokenId":          data.Id.ValueString(),
		"comment":          data.Description.ValueString(),
		"displayName":      (*string)(nil),
		"expirationInDays": (*int)(nil),
		"groups":           (*[]string)(nil),
	}

	if err := r.client.Mutate(ctx, &updateResult, variables); err == nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else {
		to_print := fmt.Sprintf("MC client 'CreateOrUpdateServiceApiToken' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
	}
}

func (r *ServiceAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServiceAccountResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteResult := client.DeleteAccessToken{}
	variables := map[string]interface{}{"tokenId": data.Id.ValueString()}

	if err := r.client.Mutate(ctx, &deleteResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'DeleteAccessToken' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
	} else if !deleteResult.DeleteAccessToken.Success {
		toPrint := "MC client 'DeleteAccessToken' mutation - success = false, " +
			"service account probably already doesn't exists. This resource will continue with its deletion"
		resp.Diagnostics.AddWarning(toPrint, "")
	}
}
