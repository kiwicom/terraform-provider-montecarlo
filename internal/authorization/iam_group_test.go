package authorization_test

import (
	"os"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/internal/acctest"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIamGroupResource(t *testing.T) {
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
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "name", "group-1"),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "label", "group-1"),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "description", ""),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "role", "mcd/editor"),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "domains.#", "0"),
					resource.TestCheckNoResourceAttr("montecarlo_iam_group.test", "ssoGroup"),
				),
			},
			{ // ImportState testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				ResourceName:      "montecarlo_iam_group.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return s.RootModule().Resources["montecarlo_iam_group.test"].Primary.Attributes["name"], nil
				},
				ImportStateVerifyIdentifierAttribute: "name",
			},
			{ // Update and Read testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigFile:               config.TestNameFile("update.tf"),
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "name", "group-1"),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "label", "group-1"),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "description", ""),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "role", "mcd/viewer"),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "domains.#", "2"),
					resource.TestCheckTypeSetElemAttr("montecarlo_iam_group.test", "domains.*", "ba0c4080-089d-4377-8878-466c31d19807"),
					resource.TestCheckTypeSetElemAttr("montecarlo_iam_group.test", "domains.*", "dd4cda19-1c5c-4339-9628-76376c9e281e"),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "sso_group", "ssoGroup1"),
				),
			},
		},
	})
}
