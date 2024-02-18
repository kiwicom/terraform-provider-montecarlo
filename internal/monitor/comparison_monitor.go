package monitor

import (
	"context"
	"fmt"

	"github.com/kiwicom/terraform-provider-montecarlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/internal/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ComparisonMonitorResource{}
var _ resource.ResourceWithImportState = &ComparisonMonitorResource{}

// To simplify provider implementations, a named function can be created with the resource implementation.
func NewComparisonMonitorResource() resource.Resource {
	return &ComparisonMonitorResource{}
}

// ComparisonMonitorResource defines the resource implementation.
type ComparisonMonitorResource struct {
	client client.MonteCarloClient
}

// ComparisonMonitorResourceModel describes the resource data model according to its Schema.
type ComparisonMonitorResourceModel struct {
	Uuid            types.String   `tfsdk:"uuid"`
	Description     types.String   `tfsdk:"description"`
	Comparisons     Comparison     `tfsdk:"comparisons"`
	QueryResultType types.String   `tfsdk:"query_result_type"`
	Source          Source         `tfsdk:"source"`
	Target          Target         `tfsdk:"target"`
	ScheduleConfig  ScheduleConfig `tfsdk:"schedule_config"`
}

type Comparison struct {
	Operator            types.String  `tfsdk:"operator"`
	ThresholdValue      types.Float64 `tfsdk:"threshold_value"`
	ComparisonType      types.String  `tfsdk:"comparison_type"`
	IsThresholdRelative types.Bool    `tfsdk:"is_threshold_relative"`
}

type Source struct {
	WarehouseUuid types.String `tfsdk:"warehouse_uuid"`
	SqlQuery      types.String `tfsdk:"sql_query"`
}

type Target struct {
	WarehouseUuid types.String `tfsdk:"warehouse_uuid"`
	SqlQuery      types.String `tfsdk:"sql_query"`
}

type ScheduleConfig struct {
	ScheduleType types.String `tfsdk:"schedule_type"`
}

type QueryResultType string

type CustomRuleComparisonInput struct {
	Operator            string  `json:"operator"`
	Threshold           float64 `json:"threshold"`
	ComparisonType      string  `json:"comparisonType"`
	IsThresholdRelative bool    `json:"isThresholdRelative"`
}

type ScheduleConfigInput struct {
	ScheduleType string `json:"scheduleType"`
}

func (r *ComparisonMonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_comparison_monitor"
}

func (r *ComparisonMonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Required:   true,
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"comparisons": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"operator": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"EQ",
								"NEQ",
								"LT",
								"LTE",
								"GT",
								"GTE",
								"IS_NULL",
								"IS_NOT_NULL",
								"AUTO",
							),
						},
					},
					"threshold_value": schema.Float64Attribute{
						Required: true,
					},
					"comparison_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"THRESHOLD",
								"DYNAMIC_THRESHOLD",
								"CHANGE",
								"FRESHNESS",
								"ABSOLUTE_VOLUME",
								"GROWTH_VOLUME",
								"QUERY_PERFORMANCE",
								"SOURCE_TARGET_DELTA",
							),
						},
					},
					"is_threshold_relative": schema.BoolAttribute{
						Computed: true,
						Optional: true,
						Default:  booldefault.StaticBool(false),
					},
				},
			},
			"query_result_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("SINGLE_NUMERIC", "ROW_COUNT", "LABELED_NUMERICS"),
				},
			},
			"source": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"warehouse_uuid": schema.StringAttribute{
						Required:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"sql_query": schema.StringAttribute{
						Required:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				},
			},
			"target": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"warehouse_uuid": schema.StringAttribute{
						Required:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"sql_query": schema.StringAttribute{
						Required:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
				},
			},
			"schedule_config": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"schedule_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("LOOSE", "FIXED", "DYNAMIC", "MANUAL"),
						},
					},
				},
			},
		},
	}
}

func (r *ComparisonMonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := common.Configure(req)
	resp.Diagnostics.Append(diags...)
	r.client = client
}

func (r *ComparisonMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ComparisonMonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResult := client.CreateOrUpdateComparisonRule{}
	scheduleConfig := ScheduleConfigInput{ScheduleType: data.ScheduleConfig.ScheduleType.ValueString()}
	comparisons := []CustomRuleComparisonInput{{
		Operator:            data.Comparisons.Operator.ValueString(),
		Threshold:           data.Comparisons.ThresholdValue.ValueFloat64(),
		ComparisonType:      data.Comparisons.ComparisonType.ValueString(),
		IsThresholdRelative: data.Comparisons.IsThresholdRelative.ValueBool(),
	}}

	variables := map[string]interface{}{
		"customRuleUuid":     (*client.UUID)(nil),
		"description":        data.Description.ValueString(),
		"comparisons":        comparisons,
		"queryResultType":    QueryResultType(data.QueryResultType.ValueString()),
		"sourceConnectionId": (*client.UUID)(nil),
		"sourceDwId":         client.UUID(data.Source.WarehouseUuid.ValueString()),
		"sourceSqlQuery":     data.Source.SqlQuery.ValueString(),
		"targetConnectionId": (*client.UUID)(nil),
		"targetDwId":         client.UUID(data.Target.WarehouseUuid.ValueString()),
		"targetSqlQuery":     data.Target.SqlQuery.ValueString(),
		"scheduleConfig":     scheduleConfig,
	}

	if err := r.client.Mutate(ctx, &createResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'createOrUpdateComparisonRule' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	data.Uuid = types.StringValue(createResult.CreateOrUpdateComparisonRule.CustomRule.Uuid)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ComparisonMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *ComparisonMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ComparisonMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *ComparisonMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
