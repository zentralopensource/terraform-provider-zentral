package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SantaScopedPathRegexResource{}
var _ resource.ResourceWithImportState = &SantaScopedPathRegexResource{}

func NewSantaScopedPathRegexResource() resource.Resource {
	return &SantaScopedPathRegexResource{}
}

// SantaScopedPathRegexResource defines the resource implementation.
type SantaScopedPathRegexResource struct {
	client *goztl.Client
}

func (r *SantaScopedPathRegexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_santa_scoped_path_regex"
}

func (r *SantaScopedPathRegexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	attributes := santaScopedConfigurationItemAttributes("path regex")
	attributes["policy"] = schema.StringAttribute{
		Description: "Policy of the pattern. Valid values are ALLOW and BLOCK. Santa evaluates the block pattern before the allow pattern, " +
			"so a path matched by both is blocked. In MONITOR mode an ALLOW entry grants nothing, and stops the executions it matches from being reported.",
		MarkdownDescription: "Policy of the pattern. Valid values are `ALLOW` and `BLOCK`. Santa evaluates the block pattern before the allow pattern, " +
			"so a path matched by both is blocked. In `MONITOR` mode an `ALLOW` entry grants nothing, and stops the executions it matches from being reported.",
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf(tfSantaPathRegexAllow, tfSantaPathRegexBlock),
		},
	}
	attributes["regex"] = schema.StringAttribute{
		Description: "ICU pattern, matched against the path of the executable, and unique in the configuration for a policy. " +
			"It is anchored at the start of the path, so it does not start with ^: use a leading .* to match anywhere. " +
			"Capture groups are not allowed, because the entries of a policy are combined into one pattern: use (?:abc). " +
			"At most 512 characters.",
		MarkdownDescription: "ICU pattern, matched against the path of the executable, and unique in the configuration for a policy. " +
			"It is anchored at the start of the path, so it does not start with `^`: use a leading `.*` to match anywhere. " +
			"Capture groups are not allowed, because the entries of a policy are combined into one pattern: use `(?:abc)`. " +
			"At most 512 characters.",
		Required: true,
		Validators: []validator.String{
			stringvalidator.LengthBetween(1, 512),
			trimmedStringValidator,
			// the server removes a leading ^, and the pattern it echoes would fail the apply
			stringvalidator.RegexMatches(regexp.MustCompile(`^[^^]`), "must not start with ^, the pattern is anchored at the start of the path"),
		},
	}

	resp.Schema = schema.Schema{
		Description: "Manages Santa scoped path regexes.",
		MarkdownDescription: "The resource `zentral_santa_scoped_path_regex` manages Santa scoped path regexes: an allowed or blocked path pattern " +
			"for some of the machines of a Santa configuration, combined with the pattern of the configuration and the other entries in scope. " +
			"The first scope level that matches a machine decides: serial numbers, then primary users, then tags. " +
			"An entry without any scope matches every machine its exclusions do not.",
		Attributes: attributes,
	}
}

func (r *SantaScopedPathRegexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SantaScopedPathRegexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data santaScopedPathRegex

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSSPR, _, err := r.client.SantaScopedPathRegexes.Create(ctx, santaScopedPathRegexRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create Santa scoped path regex, got error: %s", err),
		)
		addEndpointSupportDiagnostic(&resp.Diagnostics, err, "Santa scoped path regex", minZentralVersionSantaScopedConfigurationItems)
		return
	}

	tflog.Trace(ctx, "created a Santa scoped path regex")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaScopedPathRegexForState(ztlSSPR))...)
}

func (r *SantaScopedPathRegexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data santaScopedPathRegex

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSSPR, _, err := r.client.SantaScopedPathRegexes.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read Santa scoped path regex %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a Santa scoped path regex")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaScopedPathRegexForState(ztlSSPR))...)
}

func (r *SantaScopedPathRegexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data santaScopedPathRegex

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlSSPR, _, err := r.client.SantaScopedPathRegexes.Update(ctx, int(data.ID.ValueInt64()), santaScopedPathRegexRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update Santa scoped path regex %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a Santa scoped path regex")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, santaScopedPathRegexForState(ztlSSPR))...)
}

func (r *SantaScopedPathRegexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data santaScopedPathRegex

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SantaScopedPathRegexes.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete Santa scoped path regex %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a Santa scoped path regex")
}

func (r *SantaScopedPathRegexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "Santa scoped path regex", req, resp)
}
