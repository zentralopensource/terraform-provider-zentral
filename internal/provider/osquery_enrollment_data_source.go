package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OsqueryEnrollmentDataSource{}

func NewOsqueryEnrollmentDataSource() datasource.DataSource {
	return &OsqueryEnrollmentDataSource{}
}

// OsqueryEnrollmentDataSource defines the data source implementation.
type OsqueryEnrollmentDataSource struct {
	client *goztl.Client
}

func (d *OsqueryEnrollmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_osquery_enrollment"
}

func (d *OsqueryEnrollmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Allows details of a Osquery enrollment to be retrieved by its ID.",
		MarkdownDescription: "The data source `zentral_osquery_enrollment` allows details of a Osquery enrollment to be retrieved by its `ID`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the Osquery enrollment.",
				MarkdownDescription: "`ID` of the Osquery enrollment.",
				Required:            true,
			},
			"configuration_id": schema.Int64Attribute{
				Description:         "ID of the Osquery configuration.",
				MarkdownDescription: "`ID` of the Osquery configuration.",
				Computed:            true,
			},
			"osquery_release": schema.StringAttribute{
				Description:         "Osquery release to include in the enrollment artifacts.",
				MarkdownDescription: "Osquery release to include in the enrollment artifacts.",
				Computed:            true,
			},
			"package_url": schema.StringAttribute{
				Description:         "macOS package download URL.",
				MarkdownDescription: "macOS package download URL.",
				Computed:            true,
			},
			"script_url": schema.StringAttribute{
				Description:         "Linux script download URL.",
				MarkdownDescription: "Linux script download URL.",
				Computed:            true,
			},
			"powershell_script_url": schema.StringAttribute{
				Description:         "Powershell script download URL.",
				MarkdownDescription: "Powershell script download URL.",
				Computed:            true,
			},
			"version": schema.Int64Attribute{
				Description:         "Enrollment version.",
				MarkdownDescription: "Enrollment version.",
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				Description:         "Enrollment secret.",
				MarkdownDescription: "Enrollment secret.",
				Computed:            true,
			},
			"meta_business_unit_id": schema.Int64Attribute{
				Description:         "The ID of the meta business unit the machine will be assigned to at enrollment.",
				MarkdownDescription: "The `ID` of the meta business unit the machine will be assigned to at enrollment.",
				Computed:            true,
			},
			"tag_ids": schema.SetAttribute{
				Description:         "The IDs of the tags that the machine will get at enrollment.",
				MarkdownDescription: "The `ID`s of the tags that the machine will get at enrollment.",
				ElementType:         types.Int64Type,
				Computed:            true,
			},
			"serial_numbers": schema.SetAttribute{
				Description:         "The serial numbers the enrollment is restricted to.",
				MarkdownDescription: "The serial numbers the enrollment is restricted to.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"udids": schema.SetAttribute{
				Description:         "The UDIDs the enrollment is restricted to.",
				MarkdownDescription: "The `UDID`s the enrollment is restricted to.",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"quota": schema.Int64Attribute{
				Description:         "The number of time the enrollment can be used.",
				MarkdownDescription: "The number of time the enrollment can be used.",
				Computed:            true,
			},
		},
	}
}

func (d *OsqueryEnrollmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goztl.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *goztl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *OsqueryEnrollmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data osqueryEnrollment

	// Read Terraform enrollment data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ztlOE *goztl.OsqueryEnrollment
	var err error
	if data.ID.ValueInt64() > 0 {
		ztlOE, _, err = d.client.OsqueryEnrollments.GetByID(ctx, int(data.ID.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to get Osquery enrollment '%d' by ID, got error: %s", data.ID.ValueInt64(), err),
			)
		}
	}

	if ztlOE != nil {
		resp.Diagnostics.Append(resp.State.Set(ctx, osqueryEnrollmentForState(ztlOE))...)
	}
}
