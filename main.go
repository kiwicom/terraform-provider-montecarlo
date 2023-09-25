package main

import (
	"context"
	"flag"
	"log"

	"github.com/kiwicom/terraform-provider-monte-carlo/monte_carlo/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
)

func main() {
	flag.Parse()
	opts := providerserver.ServeOpts{Address: "registry.terraform.io/kiwicom/data-platform", Debug: false}
	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
