package resources

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	client *client.MonteCarloClient
}

// BigQueryWarehouseResourceModel describes the resource data model according to its Schema.
type BigQueryWarehouseResourceModel struct {
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
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource represents the integration of Monte Carlo with BigQuery data warehouse. " +
			"While this resource is not responsible for handling data access and other operations, such as data filtering, " +
			"it is responsible for managing the connection to BigQuery using the provided service account key.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Computed:            true,
				Optional:            false,
				MarkdownDescription: "Unique identifier of warehouse managed by this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_uuid": schema.StringAttribute{
				Computed:            true,
				Optional:            false,
				MarkdownDescription: "Unique identifier of connection responsible for communication with BigQuery.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the BigQuery warehouse as it will be presented in Monte Carlo.",
			},
			"data_collector_uuid": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Unique identifier of data collector this warehouse will be attached to. " +
					"Its not possible to change data collectors of already created warehouse, therefore if Terraform " +
					"detects change in this attribute it will plan recreation (which might not be successfull due to deletion " +
					"protection flag). Since this property is immutable in Monte Carlo warehouses it can be only changed in the configuration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"service_account_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				MarkdownDescription: "Service account key used by the warehouse connection for authentication and " +
					"authorization against BigQuery. The very same service account is used to grant required " +
					"permissions to Monte Carlo BigQuery warehouse for the data access. For more information " +
					"follow Monte Carlo documentation: https://docs.getmontecarlo.com/docs/bigquery",
			},
			"deletion_protection": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				MarkdownDescription: "Whether or not to allow Terraform to destroy the instance. Unless this field is set " +
					"to false in Terraform state, a terraform destroy or terraform apply that would delete the instance will fail. " +
					"This setting will prevent the deletion even if the real resource is already deleted.",
			},
		},
	}
}

func (r *BigQueryWarehouseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return // prevent 'nil' panic during `terraform plan`
	} else if pd, ok := req.ProviderData.(common.ProviderContext); ok {
		r.client = pd.MonteCarloClient
	} else {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected ProviderContext, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
}

func (r *BigQueryWarehouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BigQueryWarehouseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := r.addConnection(ctx, data)
	if diags.HasError() || diags.WarningsCount() > 0 {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.Uuid = result.Uuid
	data.ConnectionUuid = result.ConnectionUuid
	data.Name = result.Name
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BigQueryWarehouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BigQueryWarehouseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	type UUID string
	getResult := client.GetWarehouse{}
	variables := map[string]interface{}{"uuid": UUID(data.Uuid.ValueString())}
	query := "query getWarehouse($uuid: UUID) { getWarehouse(uuid: $uuid) { name,connections{uuid,type} } }"

	if bytes, err := r.client.ExecRaw(ctx, query, variables); err != nil && (bytes == nil || len(bytes) == 0) {
		toPrint := fmt.Sprintf("MC client 'GetWarehouse' query result - %s", err.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if jsonErr := json.Unmarshal(bytes, &getResult); jsonErr != nil {
		toPrint := fmt.Sprintf("MC client 'GetWarehouse' query failed to unmarshal data - %s", jsonErr.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if getResult.GetWarehouse == nil {
		toPrint := "MC client 'GetWarehouse' query failed to find warehouse"
		if err != nil {
			toPrint = fmt.Sprintf("%s - %s", toPrint, err.Error())
		} // response missing warehouse data may or may not contain error
		tflog.Error(ctx, toPrint)
		resp.State.RemoveResource(ctx)
		return
	}

	readConnectionUuid := types.StringNull()
	readServiceAccountKey := types.StringNull()
	for _, connection := range getResult.GetWarehouse.Connections {
		if connection.Uuid == data.ConnectionUuid.ValueString() {
			readConnectionUuid = data.ConnectionUuid
			readServiceAccountKey = data.ServiceAccountKey
		}
	}

	data.ConnectionUuid = readConnectionUuid
	data.ServiceAccountKey = readServiceAccountKey
	data.Name = types.StringValue(getResult.GetWarehouse.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BigQueryWarehouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BigQueryWarehouseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	type UUID string
	setNameResult := client.SetWarehouseName{}
	variables := map[string]interface{}{
		"dwId": UUID(data.Uuid.ValueString()),
		"name": data.Name.ValueString(),
	}

	if err := r.client.Mutate(ctx, &setNameResult, variables); err != nil {
		to_print := fmt.Sprintf("MC client 'SetWarehouseName' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(to_print, "")
		return
	}

	data.Name = types.StringValue(setNameResult.SetWarehouseName.Warehouse.Name)
	if data.ConnectionUuid.IsUnknown() || data.ConnectionUuid.IsNull() {
		if result, diags := r.addConnection(ctx, data); !diags.HasError() && diags.WarningsCount() <= 0 {
			data.ConnectionUuid = result.ConnectionUuid
		} else {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	type JSONString string
	updateResult := client.UpdateCredentials{}
	variables = map[string]interface{}{
		"changes":        JSONString(data.ServiceAccountKey.ValueString()),
		"connectionId":   UUID(data.ConnectionUuid.ValueString()),
		"shouldReplace":  true,
		"shouldValidate": false,
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

	type UUID string
	removeResult := client.RemoveConnection{}
	variables := map[string]interface{}{"connectionId": UUID(data.ConnectionUuid.ValueString())}

	if err := r.client.Mutate(ctx, &removeResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'RemoveConnection' mutation result - %s", err.Error())
		resp.Diagnostics.AddError(toPrint, "")
		return
	} else if !removeResult.RemoveConnection.Success {
		toPrint := "MC client 'RemoveConnection' mutation - success = false, " +
			"connection probably already doesnt exists. This resource will continue with its deletion"
		resp.Diagnostics.AddWarning(toPrint, "")
	}
}

func (r *BigQueryWarehouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idsImported := strings.Split(req.ID, ",")
	if len(idsImported) != 2 || idsImported[0] == "" || idsImported[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <uuid>,<connection_uuid>. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), idsImported[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("connection_uuid"), idsImported[1])...)
}

func (r *BigQueryWarehouseResource) addConnection(ctx context.Context, data BigQueryWarehouseResourceModel) (*BigQueryWarehouseResourceModel, diag.Diagnostics) {
	var diagsResult diag.Diagnostics
	type BqConnectionDetails map[string]interface{}
	testResult := client.TestBqCredentialsV2{}
	variables := map[string]interface{}{
		"validationName": "save_credentials",
		"connectionDetails": BqConnectionDetails{
			"serviceJson": b64.StdEncoding.EncodeToString(
				[]byte(data.ServiceAccountKey.ValueString()),
			),
		},
	}

	if err := r.client.Mutate(ctx, &testResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'TestBqCredentialsV2' mutation result - %s", err.Error())
		diagsResult.AddError(toPrint, "")
		return &data, diagsResult
	} else if !testResult.TestBqCredentialsV2.ValidationResult.Success {
		diags := toDiags(testResult.TestBqCredentialsV2.ValidationResult.Warnings)
		diags = append(diags, toDiags(testResult.TestBqCredentialsV2.ValidationResult.Errors)...)
		diagsResult.Append(diags...)
		return &data, diagsResult
	}

	type UUID string
	addResult := client.AddConnection{}
	var name, createWarehouseType *string = nil, nil
	warehouseUuid := data.Uuid.ValueStringPointer()
	dataCollectorUuid := data.DataCollectorUuid.ValueStringPointer()

	if warehouseUuid == nil || *warehouseUuid == "" {
		warehouseUuid = nil
		name = data.Name.ValueStringPointer()
		temp := "bigquery"
		createWarehouseType = &temp
	}

	variables = map[string]interface{}{
		"dcId":                (*UUID)(dataCollectorUuid),
		"dwId":                (*UUID)(warehouseUuid),
		"key":                 testResult.TestBqCredentialsV2.Key,
		"jobTypes":            []string{"metadata", "query_logs", "sql_query", "json_schema"},
		"name":                name,
		"connectionType":      "bigquery",
		"createWarehouseType": createWarehouseType,
	}

	if err := r.client.Mutate(ctx, &addResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'AddConnection' mutation result - %s", err.Error())
		diagsResult.AddError(toPrint, "")
		return &data, diagsResult
	}

	data.Uuid = types.StringValue(addResult.AddConnection.Connection.Warehouse.Uuid)
	data.ConnectionUuid = types.StringValue(addResult.AddConnection.Connection.Uuid)
	data.Name = types.StringValue(addResult.AddConnection.Connection.Warehouse.Name)
	return &data, diagsResult
}

func toDiags[T client.Warnings | client.Errors](in T) diag.Diagnostics {
	var diags diag.Diagnostics
	switch any(in).(type) {
	case client.Warnings:
		for _, value := range in {
			diags.AddWarning(string(value.FriendlyMessage), string(value.Resolution))
		}
	case client.Errors:
		for _, value := range in {
			diags.AddError(string(value.FriendlyMessage), string(value.Resolution))
		}
	}
	return diags
}
