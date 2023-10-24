package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/common"
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
	Description types.String   `tfsdk:"description"`
	Roles       []types.String `tfsdk:"roles"`
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
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
			},
			"roles": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf("POSTGRES", "MYSQL", "SQL-SERVER"),
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
		"roles":                normalize[string](data.Roles),
		"domainRestrictionIds": normalize[client.UUID](data.Domains),
		"ssoGroup":             data.SsoGroup.ValueStringPointer(),
	}

	if err := r.client.Mutate(ctx, &createResult, variables); err == nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else {
		to_print := fmt.Sprintf("MC client 'createOrUpdateAuthorizationGroup' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
	}
}

func (r *IamGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *IamGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *IamGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *IamGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
