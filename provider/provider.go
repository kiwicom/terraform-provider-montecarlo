package provider

import (
	"context"
	"fmt"

	"github.com/kiwicom/terraform-provider-monte-carlo/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure Provider satisfies various provider interfaces.
var _ provider.Provider = &Provider{}

type Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ProviderModel describes the provider data model.
type ProviderModel struct {
	MonteCarlo types.Object `tfsdk:"monte_carlo"`
}

type ProviderMonteCarloModel struct {
	API_KEY_ID    types.String `tfsdk:"api_key_id"`
	API_KEY_TOKEN types.String `tfsdk:"api_key_token"`
}

type ProviderContext struct {
	monteCarloClient *client.MonteCarloClient
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dataplatform"
	resp.Version = p.version
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"monte_carlo": schema.SingleNestedBlock{
				MarkdownDescription: "Monte Carlo specific inputs",
				Attributes: map[string]schema.Attribute{
					"api_key_id": schema.StringAttribute{
						MarkdownDescription: "Monte Carlo API key ID",
						Required:            true,
					},
					"api_key_token": schema.StringAttribute{
						MarkdownDescription: "Monte Carlo API key token",
						Required:            true,
					},
				},
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var monteCarlo ProviderMonteCarloModel
	data.MonteCarlo.As(ctx, &monteCarlo, basetypes.ObjectAsOptions{})
	client, err := client.NewMonteCarloClient(ctx, monteCarlo.API_KEY_ID.ValueString(), monteCarlo.API_KEY_TOKEN.ValueString())
	if err != nil {
		to_print := fmt.Sprintf("Creating MC client: %s", err.Error())
		resp.Diagnostics.AddError(to_print, "Please report this issue to the provider developers.")
		return
	}

	providerContext := ProviderContext{monteCarloClient: client}
	resp.DataSourceData = providerContext
	resp.ResourceData = providerContext
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBigQueryWarehouseResource,
	}
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{version: version}
	}
}
