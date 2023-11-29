package main

import (
	"context"
	"flag"
	"log"

	"github.com/kiwicom/terraform-provider-montecarlo/internal"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website
// Run the mocks generation tool for testing, configured by the .mockery.yaml file.
//go:generate go run github.com/vektra/mockery/v2@latest

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
)

func main() {
	flag.Parse()
	opts := providerserver.ServeOpts{Address: "registry.terraform.io/kiwicom/montecarlo", Debug: false}
	if err := providerserver.Serve(context.Background(), internal.New(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
