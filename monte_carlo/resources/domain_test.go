package resources_test

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

func TestAccDomainResource(t *testing.T) {
	domainUuid := "8bfc4"
	domainName1 := "domain1"
	domainName2 := "domain2"
	assignment1 := "MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++427a1600-2653-40c5-a1e7-5ec98703ee9d++project++gcp-project1-722af1c6"
	assignment2 := "MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++e7c59fd6-7ca8-41e7-8325-062ea38d3df5++dataset++postgre-dataset-1"

	providerContext := &common.ProviderContext{MonteCarloClient: initDomainMonteCarloClient(
		domainUuid, domainName1, domainName2, assignment1, assignment2,
	)}
	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"montecarlo": providerserver.NewProtocol6WithError(provider.New("test", providerContext)()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: domainConfig(domainName1, "Domain test description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_domain.test", "uuid", domainUuid),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "name", domainName1),
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
				ResourceName:                         "montecarlo_domain.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        domainUuid,
				ImportStateVerifyIdentifierAttribute: "uuid",
			},
			// Update and Read testing
			{
				Config: domainConfigUpdate(domainName2, "Domain test description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_domain.test", "uuid", domainUuid),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "name", domainName2),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "description", "Domain test description"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "assignments.#", "2"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "tags.#", "0"),
					resource.TestCheckTypeSetElemAttr("montecarlo_domain.test", "assignments.*", assignment1),
					resource.TestCheckTypeSetElemAttr("montecarlo_domain.test", "assignments.*", assignment2),
				),
			},
		}},
	)
}

func domainConfig(name string, description string) string {
	return fmt.Sprintf(`
provider "montecarlo" {
  account_service_key = {
    id    = "montecarlo"
  	token = "montecarlo"
  }
}

resource "montecarlo_domain" "test" {
  name        = %[1]q
  description = %[2]q
  tags        = [
	{
	  name = "owner"
	  value = "bi-internal"
	},
	{
	  name = "dataset_tables_1"
	}
  ]
}
`, name, description)
}

func domainConfigUpdate(name string, description string) string {
	return fmt.Sprintf(`
provider "montecarlo" {
  account_service_key = {
    id    = "montecarlo"
  	token = "montecarlo"
  }
}

resource "montecarlo_domain" "test" {
  name        = %[1]q
  description = %[2]q
  assignments = [
    "MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++427a1600-2653-40c5-a1e7-5ec98703ee9d++project++gcp-project1-722af1c6",
    "MCON++a84380ed-b962-4bd3-b150-04bc38a209d5++e7c59fd6-7ca8-41e7-8325-062ea38d3df5++dataset++postgre-dataset-1"
  ]
}
`, name, description)
}

func initDomainMonteCarloClient(domainUuid, domainName1, domainName2, assignment1, assignment2 string) client.MonteCarloClient {
	mcClient := cmock.MonteCarloClient{}
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.CreateOrUpdateDomain"), mock.MatchedBy(func(in map[string]interface{}) bool {
		return in["uuid"] == (*client.UUID)(nil) &&
			in["name"] == domainName1 &&
			in["description"] == "Domain test description"
	})).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.CreateOrUpdateDomain)
		arg.CreateOrUpdateDomain.Domain.Uuid = domainUuid
		arg.CreateOrUpdateDomain.Domain.Name = domainName1
		arg.CreateOrUpdateDomain.Domain.Description = "Domain test description"
		arg.CreateOrUpdateDomain.Domain.Assignments = []string{}
		arg.CreateOrUpdateDomain.Domain.Tags = []client.TagKeyValuePairOutput{
			{Name: "owner", Value: "bi-internal"},
			{Name: "dataset_tables_1"},
		}
	})

	readVariables1 := map[string]interface{}{"uuid": client.UUID(domainUuid)}
	tagsResponse := `"tags":[{"name":"owner","value":"bi-internal"},{"name":"dataset_tables_1"}]`
	readResponse1 := []byte(`{"getDomain":{"uuid":"` + domainUuid + `","name":"` + domainName1 + `","description":"Domain test description","assignments":[],` + tagsResponse + `}}`)
	mcClient.On("ExecRaw", mock.Anything, client.GetDomainQuery, readVariables1).Return(readResponse1, nil)

	deleteVariables1 := map[string]interface{}{"uuid": client.UUID(domainUuid)}
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.DeleteDomain"), deleteVariables1).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.DeleteDomain)
		arg.DeleteDomain.Deleted = 1
	})

	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.CreateOrUpdateDomain"), mock.MatchedBy(func(in map[string]interface{}) bool {
		return in["uuid"] == client.UUID(domainUuid) &&
			in["name"] == domainName2 &&
			in["description"] == "Domain test description"
	})).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.CreateOrUpdateDomain)
		arg.CreateOrUpdateDomain.Domain.Uuid = domainUuid
		arg.CreateOrUpdateDomain.Domain.Name = domainName2
		arg.CreateOrUpdateDomain.Domain.Description = "Domain test description"
		arg.CreateOrUpdateDomain.Domain.Tags = []client.TagKeyValuePairOutput{}
		arg.CreateOrUpdateDomain.Domain.Assignments = []string{assignment1, assignment2}

		mcClient.On("ExecRaw", mock.Anything, client.GetDomainQuery, readVariables1).Unset()
		assignmentsResponse := `"assignments":["` + assignment1 + `","` + assignment2 + `"]`
		readResponse2 := []byte(`{"getDomain":{"uuid":"` + domainUuid + `","name":"` + domainName2 + `","description":"Domain test description","tags":[],` + assignmentsResponse + `}}`)
		mcClient.On("ExecRaw", mock.Anything, client.GetDomainQuery, readVariables1).Return(readResponse2, nil)
	})
	return &mcClient
}
