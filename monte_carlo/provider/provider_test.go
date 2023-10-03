package provider_test

import (
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var providerFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"montecarlo": providerserver.NewProtocol6WithError(provider.New("test", nil)()),
}

func TestAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
