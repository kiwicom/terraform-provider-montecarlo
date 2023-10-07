package provider

import (
	"context"
	"fmt"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/common"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/resources"

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
	context *common.ProviderContext
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

func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "montecarlo"
	resp.Version = p.version
}

func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This open-source _Terraform_ provider enables users to seamlessly and quickly integrate the " +
			"**[Monte Carlo](https://www.montecarlodata.com/)** data reliabillity platform into their infrastructure as a code " +
			"(IaC) workflows. Provider ensures this functionality by communicating with **Monte Carlo** via its GraphQL API.",
		Attributes: map[string]schema.Attribute{
			"account_service_key": schema.SingleNestedAttribute{
				MarkdownDescription: "Monte Carlo generated **Account Service Key** used to authenticate API calls of " +
					"this provider. Should not be confused with personal API key. <br><br>For more information: " +
					"https://docs.getmontecarlo.com/docs/creating-an-api-token#creating-an-api-key <br><br>",
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Monte Carlo **Account service key** _ID_.",
						Required:            true,
						Sensitive:           true,
					},
					"token": schema.StringAttribute{
						MarkdownDescription: "Monte Carlo **Account service key** _token_.",
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

	if p.context != nil {
		resp.DataSourceData = *p.context
		resp.ResourceData = *p.context
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

	p.context = &common.ProviderContext{MonteCarloClient: client}
	resp.DataSourceData = *p.context
	resp.ResourceData = *p.context
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewBigQueryWarehouseResource,
		resources.NewDomainResource,
	}
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//NewExampleDataSource,
	}
}

func New(version string, context *common.ProviderContext) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{version: version, context: context}
	}
}
