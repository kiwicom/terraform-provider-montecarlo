package common

import (
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Cyclic types commonly shared in this provider packages

type ProviderContext struct {
	MonteCarloClient client.MonteCarloClient
}

// Monte Carlo commonly uses tags with inputs and outputs.
// For this reason lot of reasources and datasources must provide `tags` attribute as well.
// This class configures tags attribute structure for schemas using `tags` attribute.
type TagModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

func (t TagModel) toTagPair() client.TagKeyValuePairInput {
	return client.TagKeyValuePairInput{
		Name:  t.Name.ValueString(),
		Value: t.Name.ValueString(),
	}
}

func NewTagModel(in client.TagKeyValuePairOutput) TagModel {
	return TagModel{
		Name:  types.StringValue(in.Name),
		Value: types.StringValue(in.Value),
	}
}

func ToTagPairs(in []TagModel) []client.TagKeyValuePairInput {
	tagPairs := make([]client.TagKeyValuePairInput, 0, len(in))
	for _, element := range in {
		tagPairs = append(tagPairs, element.toTagPair())
	}
	return tagPairs
}

func FromTagPairs(in []client.TagKeyValuePairOutput) []TagModel {
	tagModels := make([]TagModel, 0, len(in))
	for _, element := range in {
		tagModels = append(tagModels, NewTagModel(element))
	}
	return tagModels
}
