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
	providerContext := &common.ProviderContext{MonteCarloClient: initDomainMonteCarloClient()}
	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"montecarlo": providerserver.NewProtocol6WithError(provider.New("test", providerContext)()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: domainConfig("domain1", "Domain test description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_domain.test", "uuid", "8bfc4"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "name", "domain1"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "description", "Domain test description"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "assignments.#", "0"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "tags.0.name", "owner"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "tags.0.value", "bi-internal"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "tags.1.name", "dataset_tables_1"),
					resource.TestCheckResourceAttr("montecarlo_domain.test", "tags.1.value", ""),
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

func initDomainMonteCarloClient() client.MonteCarloClient {
	mcClient := cmock.MonteCarloClient{}
	tags := []client.TagKeyValuePairInput{
		{Name: "owner", Value: "bi-internal"},
		{Name: "dataset_tables_1"},
	}
	createVariables1 := map[string]interface{}{
		"uuid":        (*client.UUID)(nil),
		"assignments": []string{},
		"tags":        tags,
		"name":        "domain1",
		"description": "Domain test description",
	}

	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.CreateOrUpdateDomain"), createVariables1).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.CreateOrUpdateDomain)
		arg.CreateOrUpdateDomain.Domain.Uuid = "8bfc4"
		arg.CreateOrUpdateDomain.Domain.Name = "domain1"
		arg.CreateOrUpdateDomain.Domain.Description = "Domain test description"
		arg.CreateOrUpdateDomain.Domain.Assignments = []string{}
		arg.CreateOrUpdateDomain.Domain.Tags = []client.TagKeyValuePairOutput{
			{Name: "owner", Value: "bi-internal"},
			{Name: "dataset_tables_1"},
		}
	})

	readVariables1 := map[string]interface{}{"uuid": client.UUID("8bfc4")}
	readQuery := "query getDomain($uuid: UUID!) { getDomain(uuid: $uuid) { uuid,name,description,tags{name,value},assignments,createdByEmail } }"
	tagsResponse := `"tags":[{"name":"owner","value":"bi-internal"},{"name":"dataset_tables_1"}]`
	readResponse1 := []byte(`{"getDomain":{"uuid":"8bfc4","name":"domain1","description":"Domain test description","assignments":[],` + tagsResponse + `}}`)
	mcClient.On("ExecRaw", mock.Anything, readQuery, readVariables1).Return(readResponse1, nil)

	deleteVariables1 := map[string]interface{}{"uuid": client.UUID("8bfc4")}
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.DeleteDomain"), deleteVariables1).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*client.DeleteDomain)
		arg.DeleteDomain.Deleted = 1
	})
	return &mcClient
}
