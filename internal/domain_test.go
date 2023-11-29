package internal_test

import (
	"os"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/internal/acctest"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDomainResource(t *testing.T) {
	mc_api_key_id := os.Getenv("MC_API_KEY_ID")
	mc_api_key_token := os.Getenv("MC_API_KEY_TOKEN")

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
					//resource.TestCheckResourceAttr("montecarlo_domain.test", "uuid", ),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "name", "domain1"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "description", "Domain test description"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "assignments.#", "0"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("montecarlo_domain.test", "tags.*", map[string]string{
						"name": "dataset_tables_1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("montecarlo_domain.test", "tags.*", map[string]string{
						"name":  "owner",
						"value": "bi-internal",
					}),
				),
			},
			{ // ImportState testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				ResourceName:      "montecarlo_domain.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return s.RootModule().Resources["montecarlo_domain.test"].Primary.Attributes["uuid"], nil
				},
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
			{ // Update and Read testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigFile:               config.TestNameFile("update.tf"),
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("montecarlo_domain.test", "uuid", domainUuid),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "name", "domain2"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "description", "Domain test description"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "assignments.#", "0"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("montecarlo_domain.test", "tags.*", map[string]string{
						"name": "dataset_tables_2",
					}),
				),
			},
			{
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigFile:               config.TestNameFile("update_assignments.tf"),
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("montecarlo_domain.test", "uuid", domainUuid),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "name", "domain2"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "description", "Domain test description 2"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "assignments.#", "0"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "tags.#", "0"),
				),
			},
		},
	})
}
