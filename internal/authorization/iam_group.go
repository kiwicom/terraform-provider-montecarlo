package authorization

import (
	"context"
	"fmt"

	"github.com/kiwicom/terraform-provider-montecarlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/internal/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IamGroupResource{}
var _ resource.ResourceWithImportState = &IamGroupResource{}

// To simplify provider implementations, a named function can be created with the resource implementation.
func NewIamGroupResource() resource.Resource {
	return &IamGroupResource{}
}

// IamGroupResource defines the resource implementation.
type IamGroupResource struct {
	client client.MonteCarloClient
}

// IamGroupResourceModel describes the resource data model according to its Schema.
type IamGroupResourceModel struct {
	Name        types.String   `tfsdk:"name"`
	Label       types.String   `tfsdk:"label"`
	Description types.String   `tfsdk:"description"`
	Role        types.String   `tfsdk:"role"`
	Domains     []types.String `tfsdk:"domains"`
	SsoGroup    types.String   `tfsdk:"sso_group"`
}

func (r *IamGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_group"
}

func (r *IamGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"label": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
			},
			"role": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"mcd/owner",
						"mcd/domains-manager",
						"mcd/responder",
						"mcd/editor",
						"mcd/viewer",
						"mcd/asset-viewer",
						"mcd/asset-editor",
					),
				},
			},
			"domains": schema.SetAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{},
					),
				),
			},
			"sso_group": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *IamGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := common.Configure(req)
	resp.Diagnostics.Append(diags...)
	r.client = client
}

func (r *IamGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IamGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResult := client.CreateOrUpdateAuthorizationGroup{}
	variables := map[string]interface{}{
		"name":                 data.Name.ValueString(),
		"label":                data.Name.ValueString(),
		"description":          data.Description.ValueString(),
		"roles":                []string{data.Role.ValueString()},
		"domainRestrictionIds": common.TfStringsTo[client.UUID](data.Domains),
		"ssoGroup":             data.SsoGroup.ValueStringPointer(),
	}

	if err := r.client.Mutate(ctx, &createResult, variables); err == nil {
		data.Label = types.StringValue(data.Name.ValueString())
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else {
		to_print := fmt.Sprintf("MC client 'createOrUpdateAuthorizationGroup' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
	}
}

func (r *IamGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IamGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResult := client.GetAuthorizationGroups{}
	variables := map[string]interface{}{}
	if err := r.client.Query(ctx, &getResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'GetAuthorizationGroups' query result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	var found *client.AuthorizationGroup
	for _, group := range getResult.GetAuthorizationGroups {
		if !group.IsManaged && group.Name == data.Name.ValueString() {
			found = &group
			break
		}
	}

	if found == nil {
		toPrint := fmt.Sprintf("MC client 'GetAuthorizationGroups' query failed to find group [name: %s]. "+
			"This resource will be removed from the Terraform state without deletion.", data.Name.ValueString())
		resp.Diagnostics.AddWarning(toPrint, "")
		resp.State.RemoveResource(ctx)
	} else {
		data.Label = types.StringValue(found.Label)
		data.Description = types.StringValue(found.Description)
		data.Role = common.TfStringsFrom(rolesToNames(found.Roles))[0]
		data.Domains = common.TfStringsFrom(domainsToUuids[string](found.DomainRestrictions))
		data.SsoGroup = types.StringPointerValue(found.SsoGroup)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *IamGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data IamGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateResult := client.CreateOrUpdateAuthorizationGroup{}
	variables := map[string]interface{}{
		"name":                 data.Name.ValueString(),
		"label":                data.Name.ValueString(),
		"description":          data.Description.ValueString(),
		"roles":                []string{data.Role.ValueString()},
		"domainRestrictionIds": common.TfStringsTo[client.UUID](data.Domains),
		"ssoGroup":             data.SsoGroup.ValueStringPointer(),
	}

	if err := r.client.Mutate(ctx, &updateResult, variables); err == nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else {
		to_print := fmt.Sprintf("MC client 'createOrUpdateAuthorizationGroup' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
	}
}

func (r *IamGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IamGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteResult := client.DeleteAuthorizationGroup{}
	variables := map[string]interface{}{"name": data.Name.ValueString()}

	if err := r.client.Mutate(ctx, &deleteResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'DeleteAuthorizationGroup' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if deleteResult.DeleteAuthorizationGroup.Deleted != 1 {
		toPrint := fmt.Sprintf("MC client 'DeleteAuthorizationGroup' mutation - deleted = %d, "+
			"expected result is 1 - more groups might have been deleted. This resource "+
			"will continue with its deletion", deleteResult.DeleteAuthorizationGroup.Deleted)
		resp.Diagnostics.AddWarning(toPrint, "")
	}
}

func (r *IamGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func rolesToNames(roles []struct{ Name string }) []string {
	result := make([]string, len(roles))
	for i, role := range roles {
		result[i] = role.Name
	}
	return result
}

func domainsToUuids[T ~string](domains []struct{ Uuid string }) []T {
	result := make([]T, len(domains))
	for i, domain := range domains {
		result[i] = T(domain.Uuid)
	}
	return result
}
