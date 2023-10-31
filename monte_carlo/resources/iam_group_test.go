package resources_test

import (
	"fmt"
	"slices"
	"strings"
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

func TestAccIamGroupResource(t *testing.T) {
	name := "group-1"
	role := "mcd/editor"
	roleUpdate := "mcd/owner"
	domainsUpdate := []string{"domain1", "domain2"}
	ssoGroupUpdate := "ssoGroup1"

	providerContext := &common.ProviderContext{MonteCarloClient: initIamGroupMonteCarloClient(
		name, role, roleUpdate, ssoGroupUpdate, domainsUpdate,
	)}
	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"montecarlo": providerserver.NewProtocol6WithError(provider.New("test", providerContext)()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: iamGroupConfig(name, role, nil, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "name", name),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "label", name),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "description", ""),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "role", role),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "domains.#", "0"),
					resource.TestCheckNoResourceAttr("montecarlo_iam_group.test", "ssoGroup"),
				),
			},
			{ // ImportState testing
				ResourceName:                         "montecarlo_iam_group.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        name,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			{ // Update and Read testing
				Config: iamGroupConfig(name, roleUpdate, &domainsUpdate, &ssoGroupUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "name", name),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "label", name),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "description", ""),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "role", roleUpdate),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "domains.#", "2"),
					resource.TestCheckTypeSetElemAttr("montecarlo_iam_group.test", "domains.*", domainsUpdate[0]),
					resource.TestCheckTypeSetElemAttr("montecarlo_iam_group.test", "domains.*", domainsUpdate[1]),
					resource.TestCheckResourceAttr("montecarlo_iam_group.test", "sso_group", ssoGroupUpdate),
				),
			},
		},
	})
}

func iamGroupConfig(name string, role string, domainRestrictions *[]string, ssoGroup *string) string {
	domainRestrictionsConfig := ""
	ssoGroupConfig := ""
	if domainRestrictions != nil {
		domainRestrictionsConfig = fmt.Sprintf("domains = %s", strings.Join(
			strings.Split(fmt.Sprintf("%q", *domainRestrictions), " "),
			", "),
		)
	}
	if ssoGroup != nil {
		ssoGroupConfig = fmt.Sprintf("sso_group = %q", *ssoGroup)
	}
	return fmt.Sprintf(`
provider "montecarlo" {
	account_service_key = {
		id    = "montecarlo"
		token = "montecarlo"
	}
}

resource "montecarlo_iam_group" "test" {
	name = %[1]q
	role = %[2]q
	%[3]s
	%[4]s
}
`, name, role, domainRestrictionsConfig, ssoGroupConfig)
}

func initIamGroupMonteCarloClient(name, role, roleUpdate, ssoGroupUpdate string, domainsUpdate []string) client.MonteCarloClient {
	mcClient := cmock.MonteCarloClient{}
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.CreateOrUpdateAuthorizationGroup"), mock.MatchedBy(func(in map[string]interface{}) bool {
		roles, rolesOk := in["roles"].([]string)
		domainRestrictions, domainRestrictionsOk := in["domainRestrictionIds"].([]client.UUID)
		return in["name"] == name &&
			in["label"] == name &&
			in["description"] == "" &&
			rolesOk && len(roles) == 1 && roles[0] == role &&
			domainRestrictionsOk && len(domainRestrictions) == 0 &&
			in["ssoGroup"] == (*string)(nil)
	})).Return(nil).Run(func(args mock.Arguments) {
		createResult := args.Get(1).(*client.CreateOrUpdateAuthorizationGroup)
		createResult.CreateOrUpdateAuthorizationGroup.AuthorizationGroup = client.AuthorizationGroup{
			Name:               name,
			Label:              name,
			Description:        "",
			Roles:              []struct{ Name string }{{Name: role}},
			DomainRestrictions: []struct{ Uuid string }{},
			SsoGroup:           (nil),
		}
	})

	mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		getResult := args.Get(1).(*client.GetAuthorizationGroups)
		getResult.GetAuthorizationGroups = []client.AuthorizationGroup{
			{
				Name:               name,
				Label:              name,
				Description:        "",
				Roles:              []struct{ Name string }{{Name: role}},
				DomainRestrictions: []struct{ Uuid string }{},
				SsoGroup:           nil,
			},
		}
	})

	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.DeleteAuthorizationGroup"), mock.MatchedBy(func(in map[string]interface{}) bool {
		return in["name"] == name
	})).Return(nil).Run(func(args mock.Arguments) {
		deleteResult := args.Get(1).(*client.DeleteAuthorizationGroup)
		deleteResult.DeleteAuthorizationGroup.Deleted = 1
	})

	// update testing
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.CreateOrUpdateAuthorizationGroup"), mock.MatchedBy(func(in map[string]interface{}) bool {
		roles, rolesOk := in["roles"].([]string)
		ssoGroup, ssoGroupOk := in["ssoGroup"].(*string)
		domainRestrictions, domainRestrictionsOk := in["domainRestrictionIds"].([]client.UUID)
		return in["name"] == name &&
			in["label"] == name &&
			in["description"] == "" &&
			rolesOk && len(roles) == 1 && roles[0] == roleUpdate &&
			domainRestrictionsOk && len(domainRestrictions) == 2 &&
			slices.Contains(domainRestrictions, client.UUID(domainsUpdate[0])) &&
			slices.Contains(domainRestrictions, client.UUID(domainsUpdate[1])) &&
			ssoGroupOk && *ssoGroup == ssoGroupUpdate
	})).Return(nil).Run(func(args mock.Arguments) {
		createResult := args.Get(1).(*client.CreateOrUpdateAuthorizationGroup)
		createResult.CreateOrUpdateAuthorizationGroup.AuthorizationGroup = client.AuthorizationGroup{
			Name:               name,
			Label:              name,
			Description:        "",
			Roles:              []struct{ Name string }{{Name: roleUpdate}},
			DomainRestrictions: []struct{ Uuid string }{{Uuid: domainsUpdate[0]}, {domainsUpdate[1]}},
			SsoGroup:           &ssoGroupUpdate,
		}

		// Read query will response with different result after update operation
		mcClient.ExpectedCalls = slices.DeleteFunc(mcClient.ExpectedCalls, func(call *mock.Call) bool {
			return call.Arguments.Is(mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything)
		})
		mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			getResult := args.Get(1).(*client.GetAuthorizationGroups)
			getResult.GetAuthorizationGroups = []client.AuthorizationGroup{
				{
					Name:               name,
					Label:              name,
					Description:        "",
					Roles:              []struct{ Name string }{{Name: roleUpdate}},
					DomainRestrictions: []struct{ Uuid string }{{Uuid: domainsUpdate[0]}, {domainsUpdate[1]}},
					SsoGroup:           &ssoGroupUpdate,
				},
			}
		})
	})
	return &mcClient
}
