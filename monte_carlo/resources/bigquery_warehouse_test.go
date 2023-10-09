package resources_test

import (
	"fmt"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"
	cmock "github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client/mock"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/common"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/provider"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBigQueryWarehouseResource(t *testing.T) {
	providerContext := &common.ProviderContext{MonteCarloClient: initBigQueryWarehouseMonteCarloClient()}
	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"montecarlo": providerserver.NewProtocol6WithError(provider.New("test", providerContext)()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: bigQueryWarehouseConfig("name1", "dataCollector1", "{}"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "uuid", "8bfc4"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "connection_uuid", "8cd5a"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "name", "name1"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "data_collector_uuid", "dataCollector1"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "service_account_key", "{}"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "deletion_protection", "false"),
				),
			},
			{ // ImportState testing
				ResourceName:                         "montecarlo_bigquery_warehouse.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "8bfc4,8cd5a,dataCollector1",
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"deletion_protection", "service_account_key"},
			},
			// Update and Read testing
			{
				Config: bigQueryWarehouseConfig("name2", "dataCollector1", "{\"json\": \"json\"}"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "uuid", "8bfc4"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "connection_uuid", "8cd5a"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "name", "name2"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "data_collector_uuid", "dataCollector1"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "service_account_key", "{\"json\": \"json\"}"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "deletion_protection", "false"),
				),
			},
		},
	})
}

func bigQueryWarehouseConfig(name string, dcid string, saKey string) string {
	return fmt.Sprintf(`
provider "montecarlo" {
  account_service_key = {
    id    = "montecarlo"
  	token = "montecarlo"
  }
}

resource "montecarlo_bigquery_warehouse" "test" {
  name                = %[1]q
  data_collector_uuid = %[2]q
  service_account_key = %[3]q
  deletion_protection = false
}
`, name, dcid, saKey)
}

func initBigQueryWarehouseMonteCarloClient() client.MonteCarloClient {
	mcClient := cmock.MonteCarloClient{}
	// Add connection operations
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.TestBqCredentialsV2"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.TestBqCredentialsV2)
		arg.TestBqCredentialsV2.Key = "testKey"
		arg.TestBqCredentialsV2.ValidationResult.Success = true
		arg.TestBqCredentialsV2.ValidationResult.Errors = client.BqTestErrors{}
		arg.TestBqCredentialsV2.ValidationResult.Warnings = client.BqTestWarnings{}
	})
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.AddConnection"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.AddConnection)
		arg.AddConnection.Connection.Uuid = "8cd5a"
		arg.AddConnection.Connection.Warehouse.Uuid = "8bfc4"
		arg.AddConnection.Connection.Warehouse.Name = "name1"
	})

	// Read operations
	readQuery := "query getWarehouse($uuid: UUID) { getWarehouse(uuid: $uuid) { name,connections{uuid,type},dataCollector{uuid} } }"
	readVariables1 := map[string]interface{}{"uuid": client.UUID("8bfc4")}
	readResponse1 := []byte(`{"getWarehouse":{"name":"name1","connections":[{"uuid":"8cd5a"}],"dataCollector":{"uuid":"dataCollector1"}}}`)
	mcClient.On("ExecRaw", mock.Anything, readQuery, readVariables1).Return(readResponse1, nil)

	// Delete operations
	deleteVariables2 := map[string]interface{}{"connectionId": client.UUID("8cd5a")}
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.RemoveConnection"), deleteVariables2).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.RemoveConnection)
		arg.RemoveConnection.Success = true
	})

	// Update operations
	updateVariables := map[string]interface{}{"dwId": client.UUID("8bfc4"), "name": "name2"}
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.SetWarehouseName"), updateVariables).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.SetWarehouseName)
		arg.SetWarehouseName.Warehouse.Uuid = "8bfc4"
		arg.SetWarehouseName.Warehouse.Name = "name2"
	})
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.UpdateCredentials"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.UpdateCredentials)
		arg.UpdateCredentials.Success = true
		// after update, read operation must return new results
		mcClient.On("ExecRaw", mock.Anything, readQuery, readVariables1).Unset()
		readResponse := []byte(`{"getWarehouse":{"name":"name2","connections":[{"uuid":"8cd5a"}],"dataCollector":{"uuid":"dataCollector1"}}}`)
		mcClient.On("ExecRaw", mock.Anything, readQuery, readVariables1).Return(readResponse, nil)
	})
	return &mcClient
}
