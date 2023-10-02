package main

import (
	"context"
	"flag"
	"log"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.16.0

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
)

func main() {
	flag.Parse()
	opts := providerserver.ServeOpts{Address: "registry.terraform.io/kiwicom/montecarlo", Debug: false}
	if err := providerserver.Serve(context.Background(), provider.New(version, nil), opts); err != nil {
		log.Fatal(err.Error())
	}
}
