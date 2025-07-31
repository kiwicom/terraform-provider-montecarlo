package warehouse

import (
	"context"
	"fmt"

	"github.com/kiwicom/terraform-provider-montecarlo/client"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WarehouseResource[T WarehouseResourceModel, K TestCredentials] interface {
	*BigQueryWarehouseResource | *TransactionalWarehouseResource
	testCredentials(ctx context.Context, data T) (K, diag.Diagnostics)
}

type WarehouseResourceModel interface {
	BigQueryWarehouseResourceModel | TransactionalWarehouseResourceModel
	GetUuid() types.String
	GetCollectorUuid() types.String
	GetName() types.String
	GetConnectionUuid() types.String
}

type TestCredentials interface {
	*client.TestBqCredentialsV2 | *client.TestDatabaseCredentials
}

func BqKeyExtractor(k *client.TestBqCredentialsV2) string {
	return k.TestBqCredentialsV2.Key
}
func TrxKeyExtractor(k *client.TestDatabaseCredentials) string {
	return k.TestDatabaseCredentials.Key
}

func addConnection[T WarehouseResource[J, K], J WarehouseResourceModel, K TestCredentials](
	ctx context.Context, mcClient client.MonteCarloClient, warehouse T, data J, connectionType string, keyExtractor func(K) string,
) (*client.AddConnection, diag.Diagnostics) {
	var diagsResult diag.Diagnostics
	testResult, credentialsDiags := warehouse.testCredentials(ctx, data)

	if testResult == nil {
		diagsResult.Append(credentialsDiags...)
		return nil, diagsResult
	}

	addResult := client.AddConnection{}
	var name, createWarehouseType *string = nil, nil
	warehouseUuid := data.GetUuid().ValueStringPointer()
	collectorUuid := data.GetCollectorUuid().ValueStringPointer()

	if warehouseUuid == nil || *warehouseUuid == "" {
		warehouseUuid = nil
		name = data.GetName().ValueStringPointer()
		temp := connectionType
		createWarehouseType = &temp
	}

	jobTypes := []string{"metadata", "sql_query"}
	if connectionType != client.TrxConnectionType {
		jobTypes = append(jobTypes, "query_logs", "json_schema")
	}

	variables := map[string]interface{}{
		"dcId":                (*client.UUID)(collectorUuid),
		"dwId":                (*client.UUID)(warehouseUuid),
		"key":                 keyExtractor(testResult),
		"jobTypes":            jobTypes,
		"name":                name,
		"connectionType":      connectionType,
		"createWarehouseType": createWarehouseType,
	}

	if err := mcClient.Mutate(ctx, &addResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'AddConnection' mutation result - %s", err.Error())
		diagsResult.AddError(toPrint, "")
		return nil, diagsResult
	} else {
		return &addResult, diagsResult
	}
}

func updateConnection[T WarehouseResource[J, K], J WarehouseResourceModel, K TestCredentials](
	ctx context.Context, mcClient client.MonteCarloClient, warehouse T, data J, keyExtractor func(K) string,
) (*client.UpdateCredentialsV2, diag.Diagnostics) {
	var diagsResult diag.Diagnostics
	testResult, credentialsDiags := warehouse.testCredentials(ctx, data)

	if testResult == nil {
		diagsResult.Append(credentialsDiags...)
		return nil, diagsResult
	}

	updateResult := client.UpdateCredentialsV2{}
	variables := map[string]interface{}{
		"connectionId":       client.UUID(data.GetConnectionUuid().ValueString()),
		"tempCredentialsKey": keyExtractor(testResult),
	}

	if err := mcClient.Mutate(ctx, &updateResult, variables); err != nil {
		toPrint := fmt.Sprintf("MC client 'UpdateCredentials' mutation result - %s", err.Error())
		diagsResult.AddError(toPrint, "")
		return nil, diagsResult
	} else if !updateResult.UpdateCredentialsV2.Success {
		toPrint := "MC client 'UpdateCredentials' mutation - success = false, " +
			"connection probably doesnt exists. Rerunning terraform operation usually helps."
		diagsResult.AddError(toPrint, "")
		return nil, diagsResult
	} else {
		return &updateResult, diagsResult
	}
}
