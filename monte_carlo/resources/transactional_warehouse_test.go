package resources_test

import (
	"fmt"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"
	cmock "github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client/mock"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/common"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccTransactionalWarehouseResource(t *testing.T) {
	name1 := "name1"
	name2 := "name2"
	uuid := "8bfc4"
	connectionUuid := "8cd5a"
	dcId := "dataCollector1"
	username1 := "user1"
	password1 := "password1"
	username2 := "user2"
	password2 := "password2"

	providerContext := &common.ProviderContext{MonteCarloClient: initTransactionalWarehouseMonteCarloClient(
		uuid, connectionUuid, dcId, name1, name2, username1, username2, password1, password2,
	)}
	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"montecarlo": providerserver.NewProtocol6WithError(provider.New("test", providerContext)()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: transactionalWarehouseConfig(name1, dcId, username1, password1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "uuid", uuid),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "connection_uuid", connectionUuid),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "name", name1),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "collector_uuid", dcId),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "deletion_protection", "false"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.host", "host"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.port", "5432"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.database", "database"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.username", username1),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.password", password1),
				),
			},
			{ // ImportState testing
				ResourceName:                         "montecarlo_transactional_warehouse.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        fmt.Sprintf("%s,%s,%s", uuid, connectionUuid, dcId),
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"db_type", "deletion_protection", "configuration"},
			},
			{ // Update and Read testing
				Config: transactionalWarehouseConfig(name2, dcId, username2, password2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "uuid", uuid),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "connection_uuid", connectionUuid),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "name", name2),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "collector_uuid", dcId),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "deletion_protection", "false"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.host", "host"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.port", "5432"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.database", "database"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.username", username2),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "configuration.password", password2),
				),
			},
		},
	})
}

func transactionalWarehouseConfig(name string, dcid string, username string, password string) string {
	return fmt.Sprintf(`
provider "montecarlo" {
  account_service_key = {
    id    = "montecarlo"
  	token = "montecarlo"
  }
}

resource "montecarlo_transactional_warehouse" "test" {
  name                = %[1]q
  collector_uuid      = %[2]q
  db_type             = "POSTGRES" # POSTGRES | MYSQL | SQL-SERVER
  deletion_protection = false

  configuration = {
    host     = "host"
    port     = 5432
    database = "database"
    username = %[3]q  #(secret)
    password = %[4]q  #(secret)
  }
}
`, name, dcid, username, password)
}

func initTransactionalWarehouseMonteCarloClient(uuid, connectionUuid, dcId, name1, name2, username1, username2, password1, password2 string) client.MonteCarloClient {
	mcClient := cmock.MonteCarloClient{}
	testKey := "testKey"
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.TestDatabaseCredentials"), mock.MatchedBy(func(in map[string]interface{}) bool {
		return in["connectionType"] == client.TransactionalConnectionType &&
			in["dbType"] == "postgres" &&
			in["host"] == "host" &&
			in["port"] == int64(5432) &&
			in["dbName"] == "database" &&
			in["user"] == username1 &&
			in["password"] == password1
	})).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.TestDatabaseCredentials)
		arg.TestDatabaseCredentials.Key = testKey
		arg.TestDatabaseCredentials.Success = true
		arg.TestDatabaseCredentials.Validations = []client.DatabaseTestDiagnostic{}
		arg.TestDatabaseCredentials.Warnings = []client.DatabaseTestDiagnostic{}
	})
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.AddConnection"), mock.MatchedBy(func(in map[string]interface{}) bool {
		return in["connectionType"] == client.TransactionalConnectionType &&
			(in["createWarehouseType"] != nil && *in["createWarehouseType"].(*string) == client.TransactionalConnectionType) &&
			(in["dcId"] != nil && *in["dcId"].(*client.UUID) == client.UUID(dcId)) &&
			(in["name"] != nil && *in["name"].(*string) == name1) &&
			in["dwId"] == (*client.UUID)(nil) && in["key"] == testKey
	})).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.AddConnection)
		arg.AddConnection.Connection.Uuid = connectionUuid
		arg.AddConnection.Connection.Warehouse.Uuid = uuid
		arg.AddConnection.Connection.Warehouse.Name = name1
	})

	// Read operations
	readVariables1 := map[string]interface{}{"uuid": client.UUID(uuid)}
	readResponse1 := []byte(fmt.Sprintf(`{"getWarehouse":{"name":%[1]q,"connections":[{"uuid":%[2]q,`+
		`"type":%[3]q}],"dataCollector":{"uuid":%[4]q}}}`, name1, connectionUuid, client.TransactionalConnectionTypeResponse, dcId))
	mcClient.On("ExecRaw", mock.Anything, client.GetWarehouseQuery, readVariables1).Return(readResponse1, nil)

	// Delete operations
	deleteVariables1 := map[string]interface{}{"connectionId": client.UUID(connectionUuid)}
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.RemoveConnection"), deleteVariables1).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.RemoveConnection)
		arg.RemoveConnection.Success = true
	})

	// Update operations
	updateVariables := map[string]interface{}{"dwId": client.UUID(uuid), "name": name2}
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.SetWarehouseName"), updateVariables).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.SetWarehouseName)
		arg.SetWarehouseName.Warehouse.Uuid = uuid
		arg.SetWarehouseName.Warehouse.Name = name2
	})
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.UpdateCredentials"), mock.MatchedBy(func(in map[string]interface{}) bool {
		return in["connectionId"] == client.UUID(connectionUuid) &&
			in["changes"] == client.JSONString(fmt.Sprintf(
				`{"db_type":%[1]q, "host": %[2]q, "port": %[3]d, "user": %[4]q, "password": %[5]q}`,
				"postgres", "host", int64(5432), username2, password2))
	})).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.UpdateCredentials)
		arg.UpdateCredentials.Success = true
		// after update, read operation must return new results
		mcClient.On("ExecRaw", mock.Anything, client.GetWarehouseQuery, readVariables1).Unset()
		readResponse := []byte(fmt.Sprintf(`{"getWarehouse":{"name":%[1]q,"connections":[{"uuid":%[2]q,`+
			`"type":%[3]q}],"dataCollector":{"uuid":%[4]q}}}`, name2, connectionUuid, client.TransactionalConnectionTypeResponse, dcId))
		mcClient.On("ExecRaw", mock.Anything, client.GetWarehouseQuery, readVariables1).Return(readResponse, nil)
	})
	return &mcClient
}
