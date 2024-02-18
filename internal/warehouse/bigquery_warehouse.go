package warehouse

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kiwicom/terraform-provider-montecarlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/internal/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &BigQueryWarehouseResource{}
var _ resource.ResourceWithImportState = &BigQueryWarehouseResource{}

// To simplify provider implementations, a named function can be created with the resource implementation.
func NewBigQueryWarehouseResource() resource.Resource {
	return &BigQueryWarehouseResource{}
}

// BigQueryWarehouseResource defines the resource implementation.
type BigQueryWarehouseResource struct {
	client client.MonteCarloClient
}

// BigQueryWarehouseResourceModel describes the resource data model according to its Schema.
type BigQueryWarehouseResourceModel struct {
	Uuid               types.String  `tfsdk:"uuid"`
	Credentials        BqCredentials `tfsdk:"credentials"`
	Name               types.String  `tfsdk:"name"`
	CollectorUuid      types.String  `tfsdk:"collector_uuid"`
	DeletionProtection types.Bool    `tfsdk:"deletion_protection"`
}

type BqCredentials struct {
	ConnectionUuid    types.String `tfsdk:"connection_uuid"`
	ServiceAccountKey types.String `tfsdk:"service_account_key"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

type BigQueryWarehouseResourceModelV1 struct {
	Uuid               types.String `tfsdk:"uuid"`
	ConnectionUuid     types.String `tfsdk:"connection_uuid"`
	Name               types.String `tfsdk:"name"`
	CollectorUuid      types.String `tfsdk:"collector_uuid"`
	ServiceAccountKey  types.String `tfsdk:"service_account_key"`
	DeletionProtection types.Bool   `tfsdk:"deletion_protection"`
}

type BigQueryWarehouseResourceModelV0 struct {
	Uuid               types.String `tfsdk:"uuid"`
	ConnectionUuid     types.String `tfsdk:"connection_uuid"`
	Name               types.String `tfsdk:"name"`
	DataCollectorUuid  types.String `tfsdk:"data_collector_uuid"`
	ServiceAccountKey  types.String `tfsdk:"service_account_key"`
	DeletionProtection types.Bool   `tfsdk:"deletion_protection"`
}

func (r *BigQueryWarehouseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bigquery_warehouse"
}

func (r *BigQueryWarehouseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Version: 2,
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"credentials": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"connection_uuid": schema.StringAttribute{
						Computed: true,
						Optional: false,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"service_account_key": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
					"updated_at": schema.StringAttribute{
						Computed: true,
						Optional: false,
					},
				},
			},
			"name": schema.StringAttribute{
				Required:   true,
				Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"collector_uuid": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"deletion_protection": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
		},
	}
}

func (r *BigQueryWarehouseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := common.Configure(req)
	resp.Diagnostics.Append(diags...)
	r.client = client
}

func (r *BigQueryWarehouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BigQueryWarehouseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := r.addConnection(ctx, data)
	resp.Diagnostics.Append(diags...)
	if result == nil {
		return
	}

	data.Uuid = result.Uuid
	data.Credentials.UpdatedAt = result.Credentials.UpdatedAt
	data.Credentials.ConnectionUuid = result.Credentials.ConnectionUuid
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BigQueryWarehouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BigQueryWarehouseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResult := client.GetWarehouse{}
	variables := map[string]interface{}{"uuid": client.UUID(data.Uuid.ValueString())}

	if bytes, err := r.client.ExecRaw(ctx, client.GetWarehouseQuery, variables); err != nil && len(bytes) == 0 {
		toPrint := fmt.Sprintf("MC client 'GetWarehouse' query result - %s", err.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if jsonErr := json.Unmarshal(bytes, &getResult); jsonErr != nil {
		toPrint := fmt.Sprintf("MC client 'GetWarehouse' query failed to unmarshal data - %s", jsonErr.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if getResult.GetWarehouse == nil {
		toPrint := fmt.Sprintf("MC client 'GetWarehouse' query failed to find warehouse [uuid: %s]. "+
			"This resource will be removed from the Terraform state without deletion.", data.Uuid.ValueString())
		if err != nil {
			toPrint = fmt.Sprintf("%s - %s", toPrint, err.Error())
		} // response missing warehouse data may or may not contain error
		resp.Diagnostics.AddWarning(toPrint, "")
		resp.State.RemoveResource(ctx)
		return
	}

	readCollectorUuid := getResult.GetWarehouse.DataCollector.Uuid
	confCollectorUuid := data.CollectorUuid.ValueString()
	if readCollectorUuid != confCollectorUuid {
		resp.Diagnostics.AddWarning(fmt.Sprintf("Obtained warehouse with [uuid: %s] but its Data "+
			"Collector UUID does not match with configured value [obtained: %s, configured: %s]. Warehouse "+
			"might have been moved to other Data Collector externally. This resource will be removed "+
			"from the Terraform state without deletion.",
			data.Uuid.ValueString(), readCollectorUuid, confCollectorUuid), "")
		resp.State.RemoveResource(ctx)
		return
	}

	readConnectionUuid := types.StringNull()
	readConnectionSAKey := types.StringNull()
	readConnectionUpdatedAt := types.StringNull()

	for _, connection := range getResult.GetWarehouse.Connections {
		if connection.Uuid == data.Credentials.ConnectionUuid.ValueString() {
			if connection.Type != client.BigQueryConnectionTypeResponse {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Obtained Warehouse [uuid: %s, connection_uuid: %s] but got unexpected connection "+
						"type '%s'.", data.Uuid.ValueString(), connection.Uuid, connection.Type),
					"Users can manually fix remote state or delete this resource from the Terraform configuration.")
				return
			}

			readConnectionUuid = data.Credentials.ConnectionUuid
			readConnectionSAKey = data.Credentials.ServiceAccountKey
			readConnectionUpdatedAt = types.StringValue(connection.UpdatedOn)
			if connection.UpdatedOn == "" {
				readConnectionUpdatedAt = types.StringValue(connection.CreatedOn)
			}
		}
	}

	if !readConnectionSAKey.IsNull() && !readConnectionUpdatedAt.Equal(data.Credentials.UpdatedAt) {
		readConnectionSAKey = types.StringValue("(unknown external value)")
	}

	data.Credentials.UpdatedAt = readConnectionUpdatedAt
	data.Credentials.ConnectionUuid = readConnectionUuid
	data.Credentials.ServiceAccountKey = readConnectionSAKey
	data.Name = types.StringValue(getResult.GetWarehouse.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BigQueryWarehouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BigQueryWarehouseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setNameResult := client.SetWarehouseName{}
	variables := map[string]interface{}{
		"dwId": client.UUID(data.Uuid.ValueString()),
		"name": data.Name.ValueString(),
	}

	if err := r.client.Mutate(ctx, &setNameResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'SetWarehouseName' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	if data.Credentials.ConnectionUuid.IsUnknown() || data.Credentials.ConnectionUuid.IsNull() {
		if result, diags := r.addConnection(ctx, data); result != nil {
			resp.Diagnostics.Append(diags...)
			data.Credentials.UpdatedAt = result.Credentials.UpdatedAt
			data.Credentials.ConnectionUuid = result.Credentials.ConnectionUuid
		} else {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	updateResult := client.UpdateCredentials{}
	variables = map[string]interface{}{
		"changes":        client.JSONString(data.Credentials.ServiceAccountKey.ValueString()),
		"connectionId":   client.UUID(data.Credentials.ConnectionUuid.ValueString()),
		"shouldReplace":  true,
		"shouldValidate": true,
	}

	if err := r.client.Mutate(ctx, &updateResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'UpdateCredentials' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if !updateResult.UpdateCredentials.Success {
		toPrint := "MC client 'UpdateCredentials' mutation - success = false, " +
			"connection probably doesnt exists. Rerunning terraform operation usually helps."
		resp.Diagnostics.AddError(toPrint, "")
		return
	}

	data.Credentials.UpdatedAt = types.StringValue(updateResult.UpdateCredentials.UpdatedAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BigQueryWarehouseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BigQueryWarehouseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.DeletionProtection.ValueBool() {
		resp.Diagnostics.AddError(
			"Failed to delete warehouse because deletion_protection is set to true. "+
				"Set it to false to proceed with warehouse deletion",
			"Deletion protection flag will prevent this resource deletion even if it was already deleted "+
				"from the real system. For reasons why this is preferred behaviour check out documentation.",
		)
		return
	}

	removeResult := client.RemoveConnection{}
	variables := map[string]interface{}{"connectionId": client.UUID(data.Credentials.ConnectionUuid.ValueString())}
	if err := r.client.Mutate(ctx, &removeResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'RemoveConnection' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if !removeResult.RemoveConnection.Success {
		toPrint := "MC client 'RemoveConnection' mutation - success = false, " +
			"connection probably already doesn't exists. This resource will continue with its deletion"
		resp.Diagnostics.AddWarning(toPrint, "")
	}
}

func (r *BigQueryWarehouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idsImported := strings.Split(req.ID, ",")
	if len(idsImported) == 3 && idsImported[0] != "" && idsImported[1] != "" && idsImported[2] != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), idsImported[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("credentials").AtName("connection_uuid"), idsImported[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collector_uuid"), idsImported[2])...)
	} else {
		resp.Diagnostics.AddError("Unexpected Import Identifier", fmt.Sprintf(
			"Expected import identifier with format: <warehouse_uuid>,<connection_uuid>,<data_collector_uuid>. Got: %q", req.ID),
		)
	}
}

func (r *BigQueryWarehouseResource) addConnection(ctx context.Context, data BigQueryWarehouseResourceModel) (*BigQueryWarehouseResourceModel, diag.Diagnostics) {
	var diagsResult diag.Diagnostics
	type BqConnectionDetails map[string]interface{}
	testResult := client.TestBqCredentialsV2{}
	variables := map[string]interface{}{
		"validationName": "save_credentials",
		"connectionDetails": BqConnectionDetails{
			"serviceJson": b64.StdEncoding.EncodeToString(
				[]byte(data.Credentials.ServiceAccountKey.ValueString()),
			),
		},
	}

	if err := r.client.Mutate(ctx, &testResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'TestBqCredentialsV2' mutation result - %s", err.Error())
		diagsResult.AddError(toPrint, "")
		return nil, diagsResult
	} else if !testResult.TestBqCredentialsV2.ValidationResult.Success {
		diags := bqTestDiagnosticToDiags(testResult.TestBqCredentialsV2.ValidationResult.Warnings)
		diags = append(diags, bqTestDiagnosticToDiags(testResult.TestBqCredentialsV2.ValidationResult.Errors)...)
		diagsResult.Append(diags...)
		return nil, diagsResult
	}

	addResult := client.AddConnection{}
	var name, createWarehouseType *string = nil, nil
	warehouseUuid := data.Uuid.ValueStringPointer()
	collectorUuid := data.CollectorUuid.ValueStringPointer()

	if warehouseUuid == nil || *warehouseUuid == "" {
		warehouseUuid = nil
		name = data.Name.ValueStringPointer()
		temp := client.BigQueryConnectionType
		createWarehouseType = &temp
	}

	variables = map[string]interface{}{
		"dcId":                (*client.UUID)(collectorUuid),
		"dwId":                (*client.UUID)(warehouseUuid),
		"key":                 testResult.TestBqCredentialsV2.Key,
		"jobTypes":            []string{"metadata", "query_logs", "sql_query", "json_schema"},
		"name":                name,
		"connectionType":      client.BigQueryConnectionType,
		"createWarehouseType": createWarehouseType,
	}

	if err := r.client.Mutate(ctx, &addResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'AddConnection' mutation result - %s", err.Error())
		diagsResult.AddError(toPrint, "")
		return nil, diagsResult
	}

	data.Uuid = types.StringValue(addResult.AddConnection.Connection.Warehouse.Uuid)
	data.Credentials.UpdatedAt = types.StringValue(addResult.AddConnection.Connection.CreatedOn)
	data.Credentials.ConnectionUuid = types.StringValue(addResult.AddConnection.Connection.Uuid)
	return &data, diagsResult
}

func bqTestDiagnosticToDiags[T client.BqTestWarnings | client.BqTestErrors](in T) diag.Diagnostics {
	var diags diag.Diagnostics
	switch any(in).(type) {
	case client.BqTestWarnings:
		for _, value := range in {
			diags.AddWarning(value.FriendlyMessage, value.Resolution)
		}
	case client.BqTestErrors:
		for _, value := range in {
			diags.AddError(value.FriendlyMessage, value.Resolution)
		}
	}
	return diags
}

func (r *BigQueryWarehouseResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Computed: true,
						Optional: false,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"connection_uuid": schema.StringAttribute{
						Computed: true,
						Optional: false,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"data_collector_uuid": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"service_account_key": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
					"deletion_protection": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorStateData BigQueryWarehouseResourceModelV0
				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)
				if !resp.Diagnostics.HasError() {
					upgradedStateData := BigQueryWarehouseResourceModelV1{
						Uuid:               priorStateData.Uuid,
						ConnectionUuid:     priorStateData.ConnectionUuid,
						CollectorUuid:      priorStateData.DataCollectorUuid,
						Name:               priorStateData.Name,
						ServiceAccountKey:  priorStateData.ServiceAccountKey,
						DeletionProtection: priorStateData.DeletionProtection,
					}
					resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
				}
			},
		},
		1: {
			PriorSchema: &schema.Schema{
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Computed: true,
						Optional: false,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"connection_uuid": schema.StringAttribute{
						Computed: true,
						Optional: false,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Required:   true,
						Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
					},
					"collector_uuid": schema.StringAttribute{
						Required: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"service_account_key": schema.StringAttribute{
						Required:  true,
						Sensitive: true,
					},
					"deletion_protection": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
				},
			},
			StateUpgrader: func(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
				var priorStateData BigQueryWarehouseResourceModelV1
				resp.Diagnostics.Append(req.State.Get(ctx, &priorStateData)...)
				if !resp.Diagnostics.HasError() {
					upgradedStateData := BigQueryWarehouseResourceModel{
						Uuid:               priorStateData.Uuid,
						CollectorUuid:      priorStateData.CollectorUuid,
						Name:               priorStateData.Name,
						DeletionProtection: priorStateData.DeletionProtection,
						Credentials: BqCredentials{
							ConnectionUuid:    priorStateData.ConnectionUuid,
							ServiceAccountKey: priorStateData.ServiceAccountKey,
							UpdatedAt:         types.StringNull(),
						},
					}
					resp.Diagnostics.Append(resp.State.Set(ctx, upgradedStateData)...)
				}
			},
		},
	}
}
