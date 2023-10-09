package resources

import (
	"context"
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PostgresWarehouseResource{}
var _ resource.ResourceWithImportState = &PostgresWarehouseResource{}

// To simplify provider implementations, a named function can be created with the resource implementation.
func NewPostgresWarehouseResource() resource.Resource {
	return &PostgresWarehouseResource{}
}

// PostgresWarehouseResource defines the resource implementation.
type PostgresWarehouseResource struct {
	client client.MonteCarloClient
}

// PostgresWarehouseResourceModel describes the resource data model according to its Schema.
type PostgresWarehouseResourceModel struct {
	Uuid               types.String  `tfsdk:"uuid"`
	ConnectionUuid     types.String  `tfsdk:"connection_uuid"`
	Name               types.String  `tfsdk:"name"`
	DataCollectorUuid  types.String  `tfsdk:"data_collector_uuid"`
	Configuration      Configuration `tfsdk:"configuration"`
	DeletionProtection types.Bool    `tfsdk:"deletion_protection"`
}

type Configuration struct {
	Host     types.String `tfsdk:"host"`
	Port     types.Int64  `tfsdk:"port"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (r *PostgresWarehouseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgres_warehouse"
}

func (r *PostgresWarehouseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This resource represents the integration of Monte Carlo with Postgres data warehouse. " +
			"While this resource is not responsible for handling data access and other operations, such as data filtering, " +
			"it is responsible for managing the connection to Postgres using the provided configuration.",
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
				MarkdownDescription: "Unique identifier of connection responsible for communication with Postgres.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the Postgre warehouse as it will be presented in Monte Carlo.",
			},
			"data_collector_uuid": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Unique identifier of data collector this warehouse will be attached to. " +
					"Its not possible to change data collectors of already created warehouses, therefore if Terraform " +
					"detects change in this attribute it will plan recreation (which might not be successfull due to deletion " +
					"protection flag). Since this property is immutable in Monte Carlo warehouses it can only be changed in the configuration",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"configuration": schema.SingleNestedAttribute{
				Required: true,
				MarkdownDescription: "Configuration used by the warehouse connection for connecting " +
					"to the Postgres database. For more information follow Monte Carlo documentation: " +
					"https://docs.getmontecarlo.com/docs/postgres",
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Database host",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"port": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Database port",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Database name",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"username": schema.StringAttribute{
						Required:            true,
						Sensitive:           true,
						MarkdownDescription: "Login username",
					},
					"password": schema.StringAttribute{
						Required:            true,
						Sensitive:           true,
						MarkdownDescription: "Login password",
					},
				},
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

func (r *PostgresWarehouseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PostgresWarehouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PostgresWarehouseResourceModel
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PostgresWarehouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PostgresWarehouseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResult := client.GetWarehouse{}
	variables := map[string]interface{}{"uuid": client.UUID(data.Uuid.ValueString())}
	query := "query getWarehouse($uuid: UUID) { getWarehouse(uuid: $uuid) { name,connections{uuid,type},dataCollector{uuid} } }"

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

	readDataCollectorUuid := getResult.GetWarehouse.DataCollector.Uuid
	confDataCollectorUuid := data.DataCollectorUuid.ValueString()
	if readDataCollectorUuid != confDataCollectorUuid {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Obtained Postgres warehouse with [uuid: %s] but its Data Collector UUID does not match with "+
				"configured value [obtained: %s, configured: %s]. Postgres warehouse might have been moved to other "+
				"Data Collector externally", data.Uuid.ValueString(), readDataCollectorUuid, confDataCollectorUuid),
			"Since its not possible for this provider to update Data Collector of Postgres warehouse, this resource "+
				"cannot continue to function properly. It is recommended to change Data Collector UUID for this "+
				"resource directly in the Terraform configuration",
		)
		return
	}

	readConnectionUuid := types.StringNull()
	readConfiguration := Configuration{
		Host:     types.StringNull(),
		Port:     types.Int64Null(),
		Name:     types.StringNull(),
		Username: types.StringNull(),
		Password: types.StringNull(),
	}

	for _, connection := range getResult.GetWarehouse.Connections {
		if connection.Uuid == data.ConnectionUuid.ValueString() {
			readConnectionUuid = data.ConnectionUuid
			readConfiguration = data.Configuration
		}
	}

	data.ConnectionUuid = readConnectionUuid
	data.Configuration = readConfiguration
	data.Name = types.StringValue(getResult.GetWarehouse.Name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PostgresWarehouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PostgresWarehouseResourceModel
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
	} else if data.ConnectionUuid.IsUnknown() || data.ConnectionUuid.IsNull() {
		if result, diags := r.addConnection(ctx, data); !diags.HasError() && diags.WarningsCount() <= 0 {
			data.ConnectionUuid = result.ConnectionUuid
		} else {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	updateResult := client.UpdateCredentials{}
	host := data.Configuration.Host.ValueString()
	port := data.Configuration.Port.ValueInt64()
	username := data.Configuration.Username.ValueString()
	password := data.Configuration.Password.ValueString()
	variables = map[string]interface{}{
		"changes":        client.JSONString(fmt.Sprintf(`{"db_type":"postgres", "host": "%s", "port": "%d", "user": "%s", "password": "%s"}`, host, port, username, password)),
		"connectionId":   client.UUID(data.ConnectionUuid.ValueString()),
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PostgresWarehouseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PostgresWarehouseResourceModel
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
	variables := map[string]interface{}{"connectionId": client.UUID(data.ConnectionUuid.ValueString())}
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

func (r *PostgresWarehouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idsImported := strings.Split(req.ID, ",")
	if len(idsImported) == 3 && idsImported[0] != "" && idsImported[1] != "" && idsImported[2] != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), idsImported[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("connection_uuid"), idsImported[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("data_collector_uuid"), idsImported[2])...)
	} else {
		resp.Diagnostics.AddError("Unexpected Import Identifier", fmt.Sprintf(
			"Expected import identifier with format: <warehouse_uuid>,<connection_uuid>,<data_collector_uuid>. Got: %q", req.ID),
		)
	}
}

func (r *PostgresWarehouseResource) addConnection(ctx context.Context, data PostgresWarehouseResourceModel) (*PostgresWarehouseResourceModel, diag.Diagnostics) {
	var diagsResult diag.Diagnostics
	testResult := client.TestDatabaseCredentials{}
	variables := map[string]interface{}{
		"connectionType": "transactional-db",
		"dbType":         "postgres",
		"host":           data.Configuration.Host.ValueString(),
		"port":           data.Configuration.Port.ValueInt64(),
		"dbName":         data.Configuration.Name.ValueString(),
		"user":           data.Configuration.Username.ValueString(),
		"password":       data.Configuration.Password.ValueString(),
	}

	if err := r.client.Mutate(ctx, &testResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'TestDatabaseCredentials' mutation result - %s", err.Error())
		diagsResult.AddError(toPrint, "")
		return &data, diagsResult
	} else if !testResult.TestDatabaseCredentials.Success {
		diags := databaseTestDiagnosticsToDiags(testResult.TestDatabaseCredentials.Warnings)
		diags = append(diags, databaseTestDiagnosticsToDiags(testResult.TestDatabaseCredentials.Validations)...)
		diagsResult.Append(diags...)
		return &data, diagsResult
	}

	addResult := client.AddConnection{}
	var name, createWarehouseType *string = nil, nil
	warehouseUuid := data.Uuid.ValueStringPointer()
	dataCollectorUuid := data.DataCollectorUuid.ValueStringPointer()

	if warehouseUuid == nil || *warehouseUuid == "" {
		warehouseUuid = nil
		name = data.Name.ValueStringPointer()
		temp := "transactional-db"
		createWarehouseType = &temp
	}

	variables = map[string]interface{}{
		"dcId":                (*client.UUID)(dataCollectorUuid),
		"dwId":                (*client.UUID)(warehouseUuid),
		"key":                 testResult.TestDatabaseCredentials.Key,
		"jobTypes":            []string{"metadata", "query_logs", "sql_query", "json_schema"},
		"name":                name,
		"connectionType":      "transactional-db",
		"createWarehouseType": createWarehouseType,
	}

	if err := r.client.Mutate(ctx, &addResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'AddConnection' mutation result - %s", err.Error())
		diagsResult.AddError(toPrint, "")
		return &data, diagsResult
	}

	data.Uuid = types.StringValue(addResult.AddConnection.Connection.Warehouse.Uuid)
	data.ConnectionUuid = types.StringValue(addResult.AddConnection.Connection.Uuid)
	return &data, diagsResult
}

func databaseTestDiagnosticsToDiags(in []client.DatabaseTestDiagnostic) diag.Diagnostics {
	var diags diag.Diagnostics
	for _, value := range in {
		diags.AddWarning(value.Message, value.Type)
	}
	return diags
}
