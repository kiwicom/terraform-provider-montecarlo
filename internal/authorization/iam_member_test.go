package authorization_test

import (
	"os"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/internal/acctest"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIamMemberResource(t *testing.T) {
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
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "group", "groups/TestAccIamMemberResource"),
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "member", "user:ndopjera@gmail.com"),
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "member_id", "21ddb883-7586-4034-9767-e5f966ec10df"),
				),
			},
			{ // ImportState testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				ResourceName:                         "montecarlo_iam_member.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        "groups/TestAccIamMemberResource,user:ndopjera@gmail.com",
				ImportStateVerifyIdentifierAttribute: "group",
			},
			{ // Update and Read testing
				ProtoV6ProviderFactories: acctest.TestAccProviderFactories,
				ConfigFile:               config.TestNameFile("update_group.tf"),
				ConfigVariables: config.Variables{
					"montecarlo_api_key_id":    config.StringVariable(mc_api_key_id),
					"montecarlo_api_key_token": config.StringVariable(mc_api_key_token),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "group", "groups/TestAccIamMemberResource2"),
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "member", "user:ndopjera@gmail.com"),
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "member_id", "21ddb883-7586-4034-9767-e5f966ec10df"),
				),
			},
		},
	})
}
