package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MonolithRepositoryResource{}
var _ resource.ResourceWithImportState = &MonolithRepositoryResource{}

func NewMonolithRepositoryResource() resource.Resource {
	return &MonolithRepositoryResource{}
}

// MonolithRepositoryResource defines the resource implementation.
type MonolithRepositoryResource struct {
	client *goztl.Client
}

func (r *MonolithRepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monolith_repository"
}

func (r *MonolithRepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages Monolith repositories.",
		MarkdownDescription: "The resource `zentral_monolith_repository` manages Monolith repositories.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the repository.",
				MarkdownDescription: "`ID` of the repository.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of the repository.",
				MarkdownDescription: "Name of the repository.",
				Required:            true,
			},
			"meta_business_unit_id": schema.Int64Attribute{
				Description:         "The ID of the meta business unit this repository is restricted to.",
				MarkdownDescription: "The `ID` of the meta business unit this repository is restricted to.",
				Optional:            true,
			},
			"backend": schema.StringAttribute{
				Description:         "Repository backend.",
				MarkdownDescription: "Repository backend.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{tfMonolithS3Backend, tfMonolithVirtualBackend}...),
				},
			},
			"s3": schema.SingleNestedAttribute{
				Description:         "S3 backend parameters.",
				MarkdownDescription: "S3 backend parameters.",
				Attributes: map[string]schema.Attribute{
					"bucket": schema.StringAttribute{
						Description:         "Name of the S3 bucket.",
						MarkdownDescription: "Name of the S3 bucket.",
						Required:            true,
					},
					"region_name": schema.StringAttribute{
						Description:         "Name of the S3 bucket region.",
						MarkdownDescription: "Name of the S3 bucket region.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"prefix": schema.StringAttribute{
						Description:         "Prefix of the Munki repository in the S3 bucket.",
						MarkdownDescription: "Prefix of the Munki repository in the S3 bucket.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"access_key_id": schema.StringAttribute{
						Description:         "AWS access key ID.",
						MarkdownDescription: "AWS access key ID.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"secret_access_key": schema.StringAttribute{
						Description:         "AWS secret access key.",
						MarkdownDescription: "AWS secret access key.",
						Sensitive:           true,
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"assume_role_arn": schema.StringAttribute{
						Description:         "ARN of the IAM role to assume.",
						MarkdownDescription: "ARN of the IAM role to assume.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"signature_version": schema.StringAttribute{
						Description:         "Version of the AWS request signature to use.",
						MarkdownDescription: "Version of the AWS request signature to use.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"endpoint_url": schema.StringAttribute{
						Description:         "S3 endpoint URL.",
						MarkdownDescription: "S3 endpoint URL.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"cloudfront_domain": schema.StringAttribute{
						Description:         "Cloudfront domain.",
						MarkdownDescription: "Cloudfront domain.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"cloudfront_key_id": schema.StringAttribute{
						Description:         "Cloudfront key ID.",
						MarkdownDescription: "Cloudfront key ID.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"cloudfront_privkey_pem": schema.StringAttribute{
						Description:         "Cloudfront private key in PEM form.",
						MarkdownDescription: "Cloudfront private key in PEM form.",
						Sensitive:           true,
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
				},
				Optional: true,
			},
		},
	}
}

func (r *MonolithRepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*goztl.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *goztl.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *MonolithRepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data monolithRepository

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMR, _, err := r.client.MonolithRepositories.Create(ctx, monolithRepositoryRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Monolith repository, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a Monolith repository")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithRepositoryForState(ztlMR))...)
}

func (r *MonolithRepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data monolithRepository

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMR, _, err := r.client.MonolithRepositories.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Monolith repository %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Monolith repository")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithRepositoryForState(ztlMR))...)
}

func (r *MonolithRepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data monolithRepository

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlMR, _, err := r.client.MonolithRepositories.Update(ctx, int(data.ID.ValueInt64()), monolithRepositoryRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Monolith repository %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Monolith repository")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, monolithRepositoryForState(ztlMR))...)
}

func (r *MonolithRepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data monolithRepository

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MonolithRepositories.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Monolith repository %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Monolith repository")
}

func (r *MonolithRepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Monolith repository", req, resp)
}
