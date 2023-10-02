package resources_test

import (
	"fmt"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/common"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/mocks"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/provider"
	"github.com/stretchr/testify/mock"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBigQueryWarehouseResource(t *testing.T) {
	providerContext := &common.ProviderContext{MonteCarloClient: initMonteCarloClient()}
	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"montecarlo": providerserver.NewProtocol6WithError(provider.New("test", providerContext)()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: basicConfig("name1", "dataCollector1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "uuid", "8bfc4"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "connection_uuid", "8cd5a"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "name", "name1"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "data_collector_uuid", "dataCollector1"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "service_account_key", "{}"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "deletion_protection", "false"),
				),
			},
		},
	})
}

func basicConfig(name string, dcid string) string {
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
  service_account_key = "{}"
  deletion_protection = false
}
`, name, dcid)
}

func initMonteCarloClient() client.MonteCarloClient {
	mcClient := mocks.MonteCarloClient{}
	// Add connection operations
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.TestBqCredentialsV2"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.TestBqCredentialsV2)
		arg.TestBqCredentialsV2.Key = "testKey"
		arg.TestBqCredentialsV2.ValidationResult.Success = true
		arg.TestBqCredentialsV2.ValidationResult.Errors = client.Errors{}
		arg.TestBqCredentialsV2.ValidationResult.Warnings = client.Warnings{}
	})
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.AddConnection"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.AddConnection)
		arg.AddConnection.Connection.Uuid = "8cd5a"
		arg.AddConnection.Connection.Warehouse.Uuid = "8bfc4"
		arg.AddConnection.Connection.Warehouse.Name = "name1"
	})

	// Read operations
	readQuery := "query getWarehouse($uuid: UUID) { getWarehouse(uuid: $uuid) { name,connections{uuid,type} } }"
	readVariables1 := map[string]interface{}{"uuid": client.UUID("8bfc4")}
	response1 := []byte(`{"getWarehouse":{"name":"name1","connections":[{"uuid":"8cd5a"}]}}`)
	mcClient.On("ExecRaw", mock.Anything, readQuery, readVariables1).Return(response1, nil)

	// Delete operations
	deleteVariables1 := map[string]interface{}{"connectionId": client.UUID("8cd5a")}
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.RemoveConnection"), deleteVariables1).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.RemoveConnection)
		arg.RemoveConnection.Success = true
	})
	return &mcClient
}
