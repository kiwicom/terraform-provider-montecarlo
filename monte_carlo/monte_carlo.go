package monte_carlo

import "github.com/kiwicom/terraform-provider-monte-carlo/monte_carlo/client"

// Cyclic types shared in this provider packages

type ProviderContext struct {
	MonteCarloClient *client.MonteCarloClient
}
