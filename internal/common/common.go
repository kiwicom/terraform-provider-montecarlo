package common

import (
	"fmt"

	"github.com/kiwicom/terraform-provider-montecarlo/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
		Value: t.Value.ValueString(),
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

func Configure[Req resource.ConfigureRequest | datasource.ConfigureRequest](req Req) (client.MonteCarloClient, diag.Diagnostics) {
	var providerData any
	switch request := any(req).(type) {
	case resource.ConfigureRequest:
		providerData = request.ProviderData
	case datasource.ConfigureRequest:
		providerData = request.ProviderData
	}

	var diags diag.Diagnostics
	if providerData == nil {
		return nil, diags // prevent 'nil' panic during `terraform plan`
	} else if pd, ok := providerData.(*ProviderContext); ok {
		return pd.MonteCarloClient, diags
	} else {
		diags.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *common.ProviderContext, got: %T. Please report this issue to the provider developers.", providerData))
		return nil, diags
	}
}

func TfStringsTo[T ~string](in []basetypes.StringValue) []T {
	res := make([]T, len(in))
	for i, element := range in {
		res[i] = T(element.ValueString())
	}
	return res
}

func TfStringsFrom(in []string) []types.String {
	res := make([]types.String, len(in))
	for i, element := range in {
		res[i] = types.StringValue(element)
	}
	return res
}
