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

// Describes the provider data model according to its Schema.
type ProviderModel struct {
	AccountServiceKey types.Object `tfsdk:"account_service_key"`
}

// Describes the provider nested object data model according to its Schema.
type ProviderAccountServiceKeyModel struct {
	ID    types.String `tfsdk:"id"`
	TOKEN types.String `tfsdk:"token"`
}

type ProviderContext struct {
	monteCarloClient *client.MonteCarloClient
}

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "monte_carlo"
	resp.Version = p.version
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"account_service_key": schema.SingleNestedBlock{
				MarkdownDescription: "Monte Carlo generated Account Service Key used to authenticate API calls of " +
					"this provider. Should not be confused with personal API key. For more information: " +
					"https://docs.getmontecarlo.com/docs/creating-an-api-token#creating-an-api-key",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Monte Carlo Account service key ID.",
						Required:            true,
						Sensitive:           true,
					},
					"token": schema.StringAttribute{
						MarkdownDescription: "Monte Carlo Account service key token.",
						Required:            true,
						Sensitive:           true,
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

	var accountServiceKey ProviderAccountServiceKeyModel
	data.AccountServiceKey.As(ctx, &accountServiceKey, basetypes.ObjectAsOptions{})
	client, err := client.NewMonteCarloClient(ctx, accountServiceKey.ID.ValueString(), accountServiceKey.TOKEN.ValueString())
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
