package common

import "github.com/kiwicom/terraform-provider-monte-carlo/monte_carlo/client"

// Cyclic types commonly shared in this provider packages

type ProviderContext struct {
	MonteCarloClient *client.MonteCarloClient
}
