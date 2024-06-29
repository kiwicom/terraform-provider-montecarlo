package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kiwicom/terraform-provider-montecarlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/internal/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DomainResource{}
var _ resource.ResourceWithImportState = &DomainResource{}

// To simplify provider implementations, a named function can be created with the resource implementation.
func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

// DomainResource defines the resource implementation.
type DomainResource struct {
	client client.MonteCarloClient
}

// DomainResourceModel describes the resource data model according to its Schema.
type DomainResourceModel struct {
	Uuid        types.String      `tfsdk:"uuid"`
	Name        types.String      `tfsdk:"name"`
	Description types.String      `tfsdk:"description"`
	Tags        []common.TagModel `tfsdk:"tags"`
	Assignments []types.String    `tfsdk:"assignments"`
}

func (r *DomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *DomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
				Optional: true,
				Default:  stringdefault.StaticString(""),
			},
			"tags": schema.SetNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Computed: true,
							Optional: true,
							Default:  stringdefault.StaticString(""),
						},
					},
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.ObjectType{},
						[]attr.Value{},
					),
				),
			},
			"assignments": schema.SetAttribute{
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
		},
	}
}

func (r *DomainResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("assignments"),
			path.MatchRoot("tags"),
		),
	}
}

func (r *DomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := common.Configure(req)
	resp.Diagnostics.Append(diags...)
	r.client = client
}

func (r *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResult := client.CreateOrUpdateDomain{}
	variables := map[string]interface{}{
		"uuid":        (*client.UUID)(nil),
		"assignments": common.TfStringsTo[string](data.Assignments),
		"tags":        common.ToTagPairs(data.Tags),
		"name":        data.Name.ValueString(),
		"description": data.Description.ValueString(),
	}

	if err := r.client.Mutate(ctx, &createResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'createOrUpdateDomain' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	data.Uuid = types.StringValue(createResult.CreateOrUpdateDomain.Domain.Uuid)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResult := client.GetDomain{}
	variables := map[string]interface{}{"uuid": client.UUID(data.Uuid.ValueString())}

	if bytes, err := r.client.ExecRaw(ctx, client.GetDomainQuery, variables); err != nil && len(bytes) == 0 {
		toPrint := fmt.Sprintf("MC client 'GetDomain' query result - %s", err.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if jsonErr := json.Unmarshal(bytes, &getResult); jsonErr != nil {
		toPrint := fmt.Sprintf("MC client 'GetDomain' query failed to unmarshal data - %s", jsonErr.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if getResult.GetDomain == nil {
		toPrint := fmt.Sprintf("MC client 'GetDomain' query failed to find domain [uuid: %s]. "+
			"This resource will be removed from the Terraform state without deletion.", data.Uuid.ValueString())
		if err != nil {
			toPrint = fmt.Sprintf("%s - %s", toPrint, err.Error())
		} // response missing domain data may or may not contain error
		resp.Diagnostics.AddWarning(toPrint, "")
		resp.State.RemoveResource(ctx)
		return
	}

	data.Tags = common.FromTagPairs(getResult.GetDomain.Tags)
	data.Assignments = common.TfStringsFrom(getResult.GetDomain.Assignments)
	data.Name = types.StringValue(getResult.GetDomain.Name)
	data.Description = types.StringValue(getResult.GetDomain.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResult := client.CreateOrUpdateDomain{}
	variables := map[string]interface{}{
		"uuid":        client.UUID(data.Uuid.ValueString()),
		"assignments": common.TfStringsTo[string](data.Assignments),
		"tags":        common.ToTagPairs(data.Tags),
		"name":        data.Name.ValueString(),
		"description": data.Description.ValueString(),
	}

	if err := r.client.Mutate(ctx, &createResult, variables); err == nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	} else {
		to_print := fmt.Sprintf("MC client 'createOrUpdateDomain' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
	}
}

func (r *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteResult := client.DeleteDomain{}
	variables := map[string]interface{}{"uuid": client.UUID(data.Uuid.ValueString())}

	if err := r.client.Mutate(ctx, &deleteResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'DeleteDomain' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(toPrint, "")
	} else if deleteResult.DeleteDomain.Deleted != 1 {
		toPrint := fmt.Sprintf("MC client 'DeleteDomain' mutation - deleted = %d, "+
			"expected result is 1 - more domains might have been deleted. This resource "+
			"will continue with its deletion", deleteResult.DeleteDomain.Deleted)
		resp.Diagnostics.AddWarning(toPrint, "")
	}
}

func (r *DomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}
