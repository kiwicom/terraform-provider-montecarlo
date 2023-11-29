package acctest

import (
	"os"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/internal"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"montecarlo": providerserver.NewProtocol6WithError(internal.New("test")()),
}

func TestAccPreCheck(t *testing.T) {
	if v := os.Getenv("MC_API_KEY_ID"); v == "" {
		t.Fatalf("'MC_API_KEY_ID' must be set for acceptance tests")
	} else if v := os.Getenv("MC_API_KEY_TOKEN"); v == "" {
		t.Fatalf("'MC_API_KEY_TOKEN' must be set for acceptance tests")
	}
}
