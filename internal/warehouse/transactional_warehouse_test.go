package warehouse_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/internal/acctest"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTransactionalWarehouseResource(t *testing.T) {
	mc_api_key_id := os.Getenv("MC_API_KEY_ID")
	mc_api_key_token := os.Getenv("MC_API_KEY_TOKEN")
	collectorUuid := "a08d23fc-00a0-4c36-b568-82e9d0e67ad8"

	pgHost := os.Getenv("PG_HOST")
	pgPortRaw := os.Getenv("PG_PORT")
	pgPort, pgPortErr := strconv.Atoi(pgPortRaw)
	pgDatabase := os.Getenv("PG_DATABASE")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")

	if pgHost == "" {
		t.Fatalf("'PG_HOST' must be set for this acceptance tests")
	} else if pgPortRaw == "" || pgPortErr != nil {
		t.Fatalf("'PG_PORT' (int) must be set for this acceptance tests")
	} else if pgDatabase == "" {
		t.Fatalf("'PG_DATABASE' must be set for this acceptance tests")
	} else if pgUser == "" {
		t.Fatalf("'PG_USER' must be set for this acceptance tests")
	} else if pgPassword == "" {
		t.Fatalf("'PG_PASSWORD' must be set for this acceptance tests")
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
					"pg_host":                  config.StringVariable(pgHost),
					"pg_port":                  config.IntegerVariable(pgPort),
					"pg_database":              config.StringVariable(pgDatabase),
					"pg_user":                  config.StringVariable(pgUser),
					"pg_password":              config.StringVariable(pgPassword),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "name", "name1"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "collector_uuid", collectorUuid),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "deletion_protection", "false"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.host", pgHost),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.port", pgPortRaw),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.database", pgDatabase),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.username", pgUser),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.password", pgPassword),
				),
			},
			{ // ImportState testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
					"pg_host":                  config.StringVariable(pgHost),
					"pg_port":                  config.IntegerVariable(pgPort),
					"pg_database":              config.StringVariable(pgDatabase),
					"pg_user":                  config.StringVariable(pgUser),
					"pg_password":              config.StringVariable(pgPassword),
				},
				ResourceName:      "montecarlo_transactional_warehouse.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					uuid := s.RootModule().Resources["montecarlo_transactional_warehouse.test"].Primary.Attributes["uuid"]
					connectionUuid := s.RootModule().Resources["montecarlo_transactional_warehouse.test"].Primary.Attributes["credentials.connection_uuid"]
					return fmt.Sprintf("%[1]s,%[2]s,%[3]s", uuid, connectionUuid, collectorUuid), nil
				},
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateVerifyIgnore:              []string{"db_type", "deletion_protection", "credentials"},
			},
			{ // Update and Read testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigFile:               config.TestNameFile("update.tf"),
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
					"pg_host":                  config.StringVariable(pgHost),
					"pg_port":                  config.IntegerVariable(pgPort),
					"pg_database":              config.StringVariable(pgDatabase),
					"pg_user":                  config.StringVariable(pgUser),
					"pg_password":              config.StringVariable(pgPassword),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "name", "name1"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "collector_uuid", collectorUuid),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "deletion_protection", "false"),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.host", pgHost),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.port", pgPortRaw),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.database", pgDatabase),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.username", pgUser),
					resource.TestCheckResourceAttr("montecarlo_transactional_warehouse.test", "credentials.password", pgPassword),
				),
			},
		},
	})
}
