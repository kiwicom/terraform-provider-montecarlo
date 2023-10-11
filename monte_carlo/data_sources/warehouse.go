package datasources

import (
	"context"
	"fmt"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/common"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &WarehouseDataSource{}

func NewWarehouseDatasource() datasource.DataSource {
	return &WarehouseDataSource{}
}

type WarehouseDataSource struct {
	client client.MonteCarloClient
}

type WarehouseDataSourceModel struct {
	Uuid     types.String                               `tfsdk:"uuid"`
	Projects map[string]WarehouseProjectDataSourceModel `tfsdk:"projects"`
}

type WarehouseProjectDataSourceModel struct {
	Mcon     types.String                               `tfsdk:"mcon"`
	Datasets map[string]WarehouseDatasetDataSourceModel `tfsdk:"datasets"`
}

type WarehouseDatasetDataSourceModel struct {
	Mcon   types.String                             `tfsdk:"mcon"`
	Tables map[string]WarehouseTableDataSourceModel `tfsdk:"tables"`
}

type WarehouseTableDataSourceModel struct {
	Mcon types.String `tfsdk:"mcon"`
}

func (d *WarehouseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse"
}

func (d *WarehouseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "",
			},
			"projects": schema.MapNestedAttribute{
				Computed:            true,
				Optional:            false,
				MarkdownDescription: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mcon": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "",
						},
						"datasets": schema.MapNestedAttribute{
							MarkdownDescription: "",
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"mcon": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "",
									},
									"tables": schema.MapNestedAttribute{
										MarkdownDescription: "",
										Required:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"mcon": schema.StringAttribute{
													Required:            true,
													MarkdownDescription: "",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *WarehouseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := common.Configure(req)
	resp.Diagnostics.Append(diags...)
	d.client = client
}

func (d *WarehouseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WarehouseDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasNextPage := true
	data.Projects = map[string]WarehouseProjectDataSourceModel{}
	variables := map[string]interface{}{
		"dwId":       client.UUID(data.Uuid.ValueString()),
		"first":      500,
		"after":      (*string)(nil),
		"isDeleted":  false,
		"isExcluded": false,
	}

	for hasNextPage {
		readResult := client.GetTables{}
		if err := d.client.Query(ctx, &readResult, variables); err != nil {
			to_print := fmt.Sprintf("MC client 'getTables' query result - %s", err.Error())
			resp.Diagnostics.AddError(to_print, "")
			return
		}

		hasNextPage = readResult.GetTables.PageInfo.HasNextPage
		variables["after"] = readResult.GetTables.PageInfo.EndCursor

		for _, element := range readResult.GetTables.Edges {
			project := data.Projects[element.Node.ProjectName]
			if project.Datasets == nil {
				project.Mcon = types.StringValue(fmt.Sprintf(
					"MCON++%s++%s++project++%s",
					element.Node.Warehouse.Account.Uuid,
					element.Node.Warehouse.Uuid,
					element.Node.ProjectName))
				project.Datasets = map[string]WarehouseDatasetDataSourceModel{}
			}

			dataset := project.Datasets[element.Node.Dataset]
			if dataset.Tables == nil {
				dataset.Mcon = types.StringValue(fmt.Sprintf(
					"MCON++%s++%s++dataset++%s:%s",
					element.Node.Warehouse.Account.Uuid,
					element.Node.Warehouse.Uuid,
					element.Node.ProjectName,
					element.Node.Dataset))
				dataset.Tables = map[string]WarehouseTableDataSourceModel{}
			}

			table := dataset.Tables[element.Node.TableId]
			table.Mcon = types.StringValue(element.Node.Mcon)
			dataset.Tables[element.Node.TableId] = table
			project.Datasets[element.Node.Dataset] = dataset
			data.Projects[element.Node.ProjectName] = project
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
