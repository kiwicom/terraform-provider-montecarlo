package datasources_test

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

func TestAccWarehouseDataSource(t *testing.T) {
	accountUuid := "a84380ed-b962-4bd3-b150-04bc38a209d5++e7c59fd6"
	warehouseUuid := "427a1600-2653-40c5-a1e7-5ec98703ee9d"
	project1 := "bi-prod"
	project2 := "booking"
	dataset1 := "raw"
	dataset2 := "processed"
	table1 := "events"
	table2 := "pageHits"
	assignment1 := fmt.Sprintf("MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++427a1600-2653-40c5-a1e7-5ec98703ee9d++table++%s:%s.%s",
		project1, dataset1, table1)
	assignment2 := fmt.Sprintf("MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++e7c59fd6-7ca8-41e7-8325-062ea38d3df5++table++%s:%s.%s",
		project2, dataset2, table2)

	providerContext := &common.ProviderContext{MonteCarloClient: initWarehouseMonteCarloClient(
		warehouseUuid, accountUuid, assignment1, assignment2, project1, project2,
		dataset1, dataset2, table1, table2,
	)}
	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"montecarlo": providerserver.NewProtocol6WithError(provider.New("test", providerContext)()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{ // Read testing
				Config: warehouseConfig(warehouseUuid),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "uuid", warehouseUuid),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects.%", "2"),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project1+".mcon",
						fmt.Sprintf("MCON++%s++%s++project++%s", accountUuid, warehouseUuid, project1)),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project2+".mcon",
						fmt.Sprintf("MCON++%s++%s++project++%s", accountUuid, warehouseUuid, project2)),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project1+".datasets."+dataset1+".mcon",
						fmt.Sprintf("MCON++%s++%s++dataset++%s:%s", accountUuid, warehouseUuid, project1, dataset1)),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project2+".datasets."+dataset2+".mcon",
						fmt.Sprintf("MCON++%s++%s++dataset++%s:%s", accountUuid, warehouseUuid, project2, dataset2)),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project1+".datasets."+dataset1+".tables."+table1+".mcon", assignment1),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project2+".datasets."+dataset2+".tables."+table2+".mcon", assignment2),
				),
			},
		},
	})
}

func warehouseConfig(uuid string) string {
	return fmt.Sprintf(`
provider "montecarlo" {
  account_service_key = {
    id    = "montecarlo"
  	token = "montecarlo"
  }
}

data "montecarlo_warehouse" "test" {
  uuid = %[1]q
}
`, uuid)
}

func initWarehouseMonteCarloClient(
	warehouseUuid, accountUuid, assignment1, assignment2,
	project1, project2, dataset1, dataset2, table1, table2 string) client.MonteCarloClient {
	mcClient := cmock.MonteCarloClient{}
	mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetTables"), mock.MatchedBy(func(in map[string]interface{}) bool {
		return in["dwId"] == client.UUID(warehouseUuid) &&
			in["isDeleted"] == false && in["isExcluded"] == false
	})).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.GetTables)
		arg.GetTables.PageInfo.HasNextPage = false

		edge1 := client.GetTablesEdge{}
		edge1.Node.Mcon = assignment1
		edge1.Node.ProjectName = project1
		edge1.Node.Dataset = dataset1
		edge1.Node.TableId = table1
		edge1.Node.Warehouse.Uuid = warehouseUuid
		edge1.Node.Warehouse.Account.Uuid = accountUuid

		edge2 := client.GetTablesEdge{}
		edge2.Node.Mcon = assignment2
		edge2.Node.ProjectName = project2
		edge2.Node.Dataset = dataset2
		edge2.Node.TableId = table2
		edge2.Node.Warehouse.Uuid = warehouseUuid
		edge2.Node.Warehouse.Account.Uuid = accountUuid

		arg.GetTables.Edges = append(arg.GetTables.Edges, edge1, edge2)
	})
	return &mcClient
}
