package warehouse_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/internal/acctest"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWarehouseDataSource(t *testing.T) {
	mc_api_key_id := os.Getenv("MC_API_KEY_ID")
	mc_api_key_token := os.Getenv("MC_API_KEY_TOKEN")

	project := "data-playground-8bb9fc23"
	dataset := "terraform_provider_montecarlo"
	personTable := "person"
	deviceTable := "device"
	accountUuid := "3e9abc75-5dc1-447e-b4cb-9d5a6fc5db5c"
	warehouseUuid := "da6c0716-2724-4bfc-b5cc-7e0364faf979"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{ // Read testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigFile:               config.TestNameFile("read.tf"),
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "uuid", "da6c0716-2724-4bfc-b5cc-7e0364faf979"),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects.%", "1"),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project+".mcon",
						fmt.Sprintf("MCON++%s++%s++project++%s", accountUuid, warehouseUuid, project)),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project+".datasets."+dataset+".mcon",
						fmt.Sprintf("MCON++%s++%s++dataset++%s:%s", accountUuid, warehouseUuid, project, dataset)),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project+".datasets."+dataset+".tables."+personTable+".mcon",
						fmt.Sprintf("MCON++%s++%s++table++%s:%s.%s", accountUuid, warehouseUuid, project, dataset, personTable)),
					resource.TestCheckResourceAttr("data.montecarlo_warehouse.test", "projects."+project+".datasets."+dataset+".tables."+deviceTable+".mcon",
						fmt.Sprintf("MCON++%s++%s++table++%s:%s.%s", accountUuid, warehouseUuid, project, dataset, deviceTable)),
				),
			},
		},
	})
}
