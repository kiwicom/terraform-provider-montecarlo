package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kiwicom/terraform-provider-montecarlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/internal/common"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MonitorResource{}
var _ resource.ResourceWithImportState = &MonitorResource{}

// To simplify provider implementations, a named function can be created with the resource implementation.
func NewMonitorResource() resource.Resource {
	return &MonitorResource{}
}

// MonitorResource defines the resource implementation.
type MonitorResource struct {
	client client.MonteCarloClient
}

// MonitorResourceModel describes the resource data model according to its Schema.
type MonitorResourceModel struct {
	Uuid     types.String `tfsdk:"uuid"`
	Resource types.String `tfsdk:"resource"`
	Monitor  types.String `tfsdk:"monitor"`
}

func (r *MonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor"
}

func (r *MonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"monitor": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *MonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := common.Configure(req)
	resp.Diagnostics.Append(diags...)
	r.client = client
}

func (r *MonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MonitorResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jsonMonitor, err := yamlToJson(data.Monitor.ValueString())
	if err != nil {
		summary := "Failed to convert Monitor YAML definition to JSON: " + err.Error()
		detail := "Conversion must be done before sending the request to the API"
		resp.Diagnostics.AddError(summary, detail)
		return
	}

	namespace := uuid.New().String()
	createResult := client.CreateOrUpdateMonteCarloConfigTemplate{}
	variables := map[string]interface{}{
		"namespace":          namespace,
		"resource":           data.Resource.ValueStringPointer(),
		"dryRun":             true,
		"configTemplateJson": jsonMonitor,
	}

	response := &createResult.CreateOrUpdateMonteCarloConfigTemplate.Response
	if err := r.client.Mutate(ctx, &createResult, variables); err != nil {
		summary := fmt.Sprintf("MC client 'CreateOrUpdateMonteCarloConfigTemplate' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(summary, "")
		return
	} else if !response.ChangesApplied || response.ErrorsAsJson == "" {
		summary := fmt.Sprintf("MC client 'CreateOrUpdateMonteCarloConfigTemplate' mutation result - %s", response.ErrorsAsJson)
		resp.Diagnostics.AddError(summary, "")
		return
	}

	data.Uuid = types.StringValue(namespace)
	data.Resource = types.StringValue("")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *MonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *MonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *MonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}

func yamlToJson(yamlData string) (string, error) {
	var data map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlData), &data); err != nil {
		return "", err
	}
	jsonData, err := json.Marshal(data)
	return string(jsonData), err
}
