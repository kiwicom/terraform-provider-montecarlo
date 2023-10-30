package resources

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var groupsRegex = regexp.MustCompile(`^groups/.+$`)
var memberRegex = regexp.MustCompile(`^user/.+$`)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &IamMemberResource{}
var _ resource.ResourceWithImportState = &IamMemberResource{}

// To simplify provider implementations, a named function can be created with the resource implementation.
func NewIamMemberResource() resource.Resource {
	return &IamMemberResource{}
}

// IamMemberResource defines the resource implementation.
type IamMemberResource struct {
	client client.MonteCarloClient
}

// IamMemberResourceModel describes the resource data model according to its Schema.
type IamMemberResourceModel struct {
	Group    types.String `tfsdk:"group"`
	Member   types.String `tfsdk:"member"`
	MemberId types.String `tfsdk:"member_id"`
}

func (r *IamMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_member"
}

func (r *IamMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"group": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(groupsRegex, "Expected format: groups/{group_name}"),
				},
			},
			"member": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(memberRegex, "Expected format: user/{user_email}"),
				},
			},
			"member_id": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *IamMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := common.Configure(req)
	resp.Diagnostics.Append(diags...)
	r.client = client
}

func (r *IamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data IamMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userEmail := strings.Split(data.Member.ValueString(), "user/")[1]
	getUserResult := client.GetUsersInAccount{}
	variables := map[string]interface{}{
		"email": userEmail,
		"first": 1,
		"after": (*string)(nil),
	}

	if err := r.client.Query(ctx, &getUserResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'getTables' query result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	} else if len(getUserResult.GetUsersInAccount.Edges) == 0 {
		to_print := fmt.Sprintf("User %s not found", userEmail)
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	groupName := strings.Split(data.Group.ValueString(), "groups/")[1]
	getGroupResult := client.GetAuthorizationGroups{}
	variables = map[string]interface{}{}
	if err := r.client.Query(ctx, &getGroupResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'GetAuthorizationGroups' query result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	var found *client.AuthorizationGroup
	if index := slices.IndexFunc(getGroupResult.GetAuthorizationGroups, func(group client.AuthorizationGroup) bool {
		return !group.IsManaged && group.Name == groupName
	}); index >= 0 && getGroupResult.GetAuthorizationGroups[index].SsoGroup == nil {
		found = &getGroupResult.GetAuthorizationGroups[index]
	} else {
		to_print := fmt.Sprintf("Group %s not found or is SSO managed", data.Group.ValueString())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	memberUserIds := make([]string, len(found.Users)+1)
	memberUserIds[len(found.Users)] = getUserResult.GetUsersInAccount.Edges[0].Node.Id
	for i, user := range found.Users {
		memberUserIds[i] = user.Id
	}

	updateResult := client.CreateOrUpdateAuthorizationGroup{}
	variables = map[string]interface{}{
		"name":                 found.Name,
		"label":                found.Label,
		"description":          found.Description,
		"roles":                rolesToNames(found.Roles),
		"domainRestrictionIds": domainsToUuids(found.DomainRestrictions),
		"ssoGroup":             found.SsoGroup,
		"memberUserIds":        memberUserIds,
	}

	if err := r.client.Mutate(ctx, &updateResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'createOrUpdateAuthorizationGroup' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
	} else {
		data.MemberId = types.StringValue(getUserResult.GetUsersInAccount.Edges[0].Node.Id)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *IamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data IamMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userEmail := strings.Split(data.Member.ValueString(), "user/")[1]
	getUserResult := client.GetUsersInAccount{}
	variables := map[string]interface{}{
		"email": userEmail,
		"first": 1,
		"after": (*string)(nil),
	}

	if err := r.client.Query(ctx, &getUserResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'getTables' query result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	} else if len(getUserResult.GetUsersInAccount.Edges) == 0 {
		to_print := fmt.Sprintf("User %s not found", userEmail)
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	groupName := strings.Split(data.Group.ValueString(), "groups/")[1]
	getGroupResult := client.GetAuthorizationGroups{}
	variables = map[string]interface{}{}
	if err := r.client.Query(ctx, &getGroupResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'GetAuthorizationGroups' query result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	var found *client.AuthorizationGroup
	for _, group := range getGroupResult.GetAuthorizationGroups {
		if !group.IsManaged && group.Name == groupName {
			found = &group
			break
		}
	}

	if found == nil || found.SsoGroup != nil {
		to_print := fmt.Sprintf("Group %s not found or is SSO managed", data.Group.ValueString())
		resp.Diagnostics.AddError(to_print, "")
	} else if !slices.Contains(found.Users, getUserResult.GetUsersInAccount.Edges[0].Node) {
		to_print := fmt.Sprintf("User %s not found in group %s", userEmail, data.Group.ValueString())
		resp.Diagnostics.AddWarning(to_print, "")
		resp.State.RemoveResource(ctx)
	} else {
		data.MemberId = types.StringValue(getUserResult.GetUsersInAccount.Edges[0].Node.Id)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func (r *IamMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO
}

func (r *IamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data IamMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupName := strings.Split(data.Group.ValueString(), "groups/")[1]
	getGroupResult := client.GetAuthorizationGroups{}
	variables := map[string]interface{}{}
	if err := r.client.Query(ctx, &getGroupResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'GetAuthorizationGroups' query result - %s", err.Error())
		resp.Diagnostics.AddWarning(to_print, "")
		return
	}

	var found *client.AuthorizationGroup
	if index := slices.IndexFunc(getGroupResult.GetAuthorizationGroups, func(group client.AuthorizationGroup) bool {
		return !group.IsManaged && group.Name == groupName
	}); index >= 0 && getGroupResult.GetAuthorizationGroups[index].SsoGroup == nil {
		found = &getGroupResult.GetAuthorizationGroups[index]
	} else {
		to_print := fmt.Sprintf("Group %s not found or is SSO managed", data.Group.ValueString())
		resp.Diagnostics.AddWarning(to_print, "")
		return
	}

	memberUserIds := make([]string, len(found.Users))
	for i, user := range found.Users {
		memberUserIds[i] = user.Id
	}

	updateResult := client.CreateOrUpdateAuthorizationGroup{}
	memberUserIds = slices.DeleteFunc(memberUserIds, func(userId string) bool { return userId == data.MemberId.ValueString() })
	variables = map[string]interface{}{
		"name":                 found.Name,
		"label":                found.Label,
		"description":          found.Description,
		"roles":                rolesToNames(found.Roles),
		"domainRestrictionIds": domainsToUuids(found.DomainRestrictions),
		"ssoGroup":             found.SsoGroup,
		"memberUserIds":        memberUserIds,
	}

	if err := r.client.Mutate(ctx, &updateResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'createOrUpdateAuthorizationGroup' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
	}
}

func (r *IamMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// TODO
}
