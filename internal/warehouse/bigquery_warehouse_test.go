package warehouse_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/internal/acctest"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBigQueryWarehouseResource(t *testing.T) {
	mc_api_key_id := os.Getenv("MC_API_KEY_ID")
	mc_api_key_token := os.Getenv("MC_API_KEY_TOKEN")

	collectorUuid := "a08d23fc-00a0-4c36-b568-82e9d0e67ad8"
	createSa, createSaErr := os.ReadFile("testdata/TestAccBigQueryWarehouseResource/create-sa.json")
	updateSa, updateSaErr := os.ReadFile("testdata/TestAccBigQueryWarehouseResource/update-sa.json")

	if createSaErr != nil {
		t.Fatalf("failed to read testdata/TestAccBigQueryWarehouseResource/create-sa.json: %v", createSaErr)
	} else if updateSaErr != nil {
		t.Fatalf("failed to read testdata/TestAccBigQueryWarehouseResource/update-sa.json: %v", updateSaErr)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acctest.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{ // Create and Read testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigFile:               config.TestNameFile("create.tf"),
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "name", "test-warehouse"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "collector_uuid", collectorUuid),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "service_account_key", string(createSa)),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "deletion_protection", "false"),
				),
			},
			{ // ImportState testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				ResourceName:      "montecarlo_bigquery_warehouse.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					uuid := s.RootModule().Resources["montecarlo_bigquery_warehouse.test"].Primary.Attributes["uuid"]
					connectionUuid := s.RootModule().Resources["montecarlo_bigquery_warehouse.test"].Primary.Attributes["connection_uuid"]
					return fmt.Sprintf("%[1]s,%[2]s,%[3]s", uuid, connectionUuid, collectorUuid), nil
				},
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"deletion_protection", "service_account_key"},
			},
			{ // Update and Read testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigFile:               config.TestNameFile("update.tf"),
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "name", "test-warehouse-updated"),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "collector_uuid", collectorUuid),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "service_account_key", string(updateSa)),
					resource.TestCheckResourceAttr("montecarlo_bigquery_warehouse.test", "deletion_protection", "false"),
				),
			},
		},
	})
}
