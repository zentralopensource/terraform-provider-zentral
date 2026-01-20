package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zentralopensource/goztl"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &MDMDEPEnrollmentResource{}
var _ resource.ResourceWithImportState = &MDMDEPEnrollmentResource{}

func NewMDMDEPEnrollmentResource() resource.Resource {
	return &MDMDEPEnrollmentResource{}
}

// MDMDEPEnrollmentResource defines the resource implementation.
type MDMDEPEnrollmentResource struct {
	client *goztl.Client
}

func (r *MDMDEPEnrollmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_dep_enrollment"
}

var mdeAuthenticationSchema schema.SingleNestedAttribute = schema.SingleNestedAttribute{
	Description:         "ADE/DEP authentication.",
	MarkdownDescription: "ADE/DEP authentication.",
	Attributes: map[string]schema.Attribute{
		"realm_uuid": schema.StringAttribute{
			Description:         "UUID of the realm used to authenticate the users during ADE/DEP enrollment.",
			MarkdownDescription: "`UUID` of the realm used to authenticate the users during ADE/DEP enrollment.",
			Required:            true,
		},
		"use_for_setup_assistant_user": schema.BoolAttribute{
			Description:         "If true, the realm user details are used to prefill the Setup Assistant user form.",
			MarkdownDescription: "If `true`, the realm user details are used to prefill the Setup Assistant user form.",
			Required:            true,
		},
		"setup_assistant_user_is_admin": schema.BoolAttribute{
			Description:         "If true, the user created with the Setup Assistant form will be admin.",
			MarkdownDescription: "If `true`, the user created with the Setup Assistant form will be admin.",
			Required:            true,
		},
		"setup_assistant_username_pattern": schema.StringAttribute{
			Description:         "Pattern used to derive the username from the realm user details. Either  $REALM_USER.DEVICE_USERNAME or $REALM_USER.EMAIL_PREFIX.",
			MarkdownDescription: "Pattern used to derive the username from the realm user details. Either  `$REALM_USER.DEVICE_USERNAME` or `$REALM_USER.EMAIL_PREFIX`.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf([]string{"$REALM_USER.DEVICE_USERNAME", "$REALM_USER.EMAIL_PREFIX"}...),
			},
		},
	},
	Optional: true,
	Computed: true,
	Default:  objectdefault.StaticValue(mdeAuthenticationDefault()),
}

var mdeExtraAdminSchema schema.SingleNestedAttribute = schema.SingleNestedAttribute{
	Description:         "Extra admin created automatically during Setup Assistant.",
	MarkdownDescription: "Extra admin created automatically during Setup Assistant.",
	Attributes: map[string]schema.Attribute{
		"hidden": schema.BoolAttribute{
			Description:         "Toggles if the created admin is hidden. Default is false.",
			MarkdownDescription: "Toggles if the created admin is hidden. Default is `false`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"full_name": schema.StringAttribute{
			Description:         "Admin full name.",
			MarkdownDescription: "Admin full name.",
			Required:            true,
		},
		"short_name": schema.StringAttribute{
			Description:         "Admin short name.",
			MarkdownDescription: "Admin short name.",
			Required:            true,
		},
		"password_complexity": schema.Int64Attribute{
			Description:         "Complexity (1 → 3) of the device specific random password. Default is 3.",
			MarkdownDescription: "Complexity (`1` → `3`) of the device specific random password. Default is `3`.",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(tfMDMDEPEnrollmentDefaultAdminPasswordComplexity),
			Validators: []validator.Int64{
				int64validator.Between(
					tfMDMDEPEnrollmentMinAdminPasswordComplexity,
					tfMDMDEPEnrollmentMaxAdminPasswordComplexity,
				),
			},
		},
		"password_rotation_delay": schema.Int64Attribute{
			Description:         "Delay in minutes (0 → 1440, 1 day) after which an automatic password rotation is triggered when the password is revealed. Default is 1 hour.",
			MarkdownDescription: "Delay in minutes (`0` → `1440`, 1 day) after which an automatic password rotation is triggered when the password is revealed. Default is 1 hour.",
			Optional:            true,
			Computed:            true,
			Default:             int64default.StaticInt64(tfMDMDEPEnrollmentDefaultAdminPasswordRotationDelay),
			Validators: []validator.Int64{
				int64validator.Between(
					tfMDMDEPEnrollmentMinAdminPasswordRotationDelay,
					tfMDMDEPEnrollmentMaxAdminPasswordRotationDelay,
				),
			},
		},
	},
	Optional: true,
	Computed: true,
	Default:  objectdefault.StaticValue(mdeExtraAdminDefault()),
}

var mdeProfileSchema schema.SingleNestedAttribute = schema.SingleNestedAttribute{
	Description:         "ADE/DEP profile settings.",
	MarkdownDescription: "ADE/DEP profile settings. See Apple [docs](https://developer.apple.com/documentation/devicemanagement/profile/).",
	Attributes: map[string]schema.Attribute{
		"virtual_server_id": schema.Int64Attribute{
			Description:         "ID of the DEP virtual server (ABM/ASM) to push the profile to.",
			MarkdownDescription: "`ID` of the DEP virtual server (ABM/ASM) to push the profile to.",
			Required:            true,
		},
		"name": schema.StringAttribute{
			Description:         "Name of the profile in ABM/ASM.",
			MarkdownDescription: "Name of the profile in ABM/ASM.",
			Required:            true,
		},
		"allow_pairing": schema.BoolAttribute{
			Description:         "Allow the device to connect to other computers. In iOS 13, this property was deprecated. Default is true.",
			MarkdownDescription: "Allow the device to connect to other computers. In iOS 13, this property was deprecated. Default is `true`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"auto_advance_setup": schema.BoolAttribute{
			Description:         "If set to true, the device will tell Setup Assistant to automatically advance though its screens. Default is false.",
			MarkdownDescription: "If set to `true`, the device will tell Setup Assistant to automatically advance though its screens. Default is `false`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"await_device_configured": schema.BoolAttribute{
			Description:         "If true, the device will not continue in Setup Assistant until the MDM server sends a command that states the device is configured. Required when using the authenticated user details or for the extra admin. Default is false.",
			MarkdownDescription: "If `true`, the device will not continue in Setup Assistant until the MDM server sends a command that states the device is configured. Required when using the authenticated user details or for the extra admin. Default is `false`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"is_mandatory": schema.BoolAttribute{
			Description:         "If true, the user may not skip applying the profile returned by the MDM server. Default is true.",
			MarkdownDescription: "If `true`, the user may not skip applying the profile returned by the MDM server. Default is `true`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"is_mdm_removable": schema.BoolAttribute{
			Description:         "If false, the MDM payload cannot be removed by the user via the user interface on the device. This key can be set to false only if is_supervised is set to true. Default is false.",
			MarkdownDescription: "If `false`, the MDM payload cannot be removed by the user via the user interface on the device. This key can be set to false only if `is_supervised` is set to `true`. Default is `false`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"is_multi_user": schema.BoolAttribute{
			Description:         "If true, tells the device to configure for Shared iPad. Default is false.",
			MarkdownDescription: "If `true`, tells the device to configure for Shared iPad. Default is `false`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"is_supervised": schema.BoolAttribute{
			Description:         "If true, the device must be supervised. Default is true.",
			MarkdownDescription: "If `true`, the device must be supervised. Default is `true`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"include_anchor_certs": schema.BoolAttribute{
			Description:         "If true, only the configured TLS certificates for the Zentral instance will be used by the device when evaluating the trust of the connection to the MDM server URL. Otherwise, the device uses the built-in root certificates. Default is false.",
			MarkdownDescription: "If `true`, only the configured TLS certificates for the Zentral instance will be used by the device when evaluating the trust of the connection to the MDM server URL. Otherwise, the device uses the built-in root certificates. Default is `false`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},

		"skip_setup_items": schema.SetAttribute{
			Description:         "A list of setup panes to skip.",
			MarkdownDescription: "A list of setup panes to skip. The list of valid strings is defined in [SkipKeys](https://developer.apple.com/documentation/devicemanagement/skipkeys).",
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
			Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
		},
		"language": schema.StringAttribute{
			Description:         "Two-letter ISO 639-1 code of the enrollment language.",
			MarkdownDescription: "Two-letter ISO 639-1 code of the enrollment language.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
		"region": schema.StringAttribute{
			Description:         "Two-letter ISO 3166-1 country code.",
			MarkdownDescription: "Two-letter ISO 3166-1 country code.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
		"department": schema.StringAttribute{
			Description:         "The user-defined department or location name.",
			MarkdownDescription: "The user-defined department or location name.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
		"org_magic": schema.StringAttribute{
			Description:         "A string that uniquely identifies various services that are managed by a single organization.",
			MarkdownDescription: "A string that uniquely identifies various services that are managed by a single organization.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
		"support_email_address": schema.StringAttribute{
			Description:         "A support email address for the organization.",
			MarkdownDescription: "A support email address for the organization.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
		"support_phone_number": schema.StringAttribute{
			Description:         "A support phone number for the organization.",
			MarkdownDescription: "A support phone number for the organization.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
	},
	Required: true,
}

var mdeOSVersionEnforcementSchema schema.SingleNestedAttribute = schema.SingleNestedAttribute{
	Description:         "OS version enforcement settings.",
	MarkdownDescription: "OS version enforcement settings.",
	Attributes: map[string]schema.Attribute{
		"ios_min_version": schema.StringAttribute{
			Description:         "The fixed minimum iOS version required for a successful enrollment.",
			MarkdownDescription: "The fixed minimum iOS version required for a successful enrollment.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
		"auto_ios_min_version_until": schema.StringAttribute{
			Description:         "If set, the minimum iOS version required for a successful enrollment will be the latest available for the enrolling device until this version (excluded). Set for example 28 to automatically required the latest iOS version until (but not including) iOS 28.",
			MarkdownDescription: "If set, the minimum iOS version required for a successful enrollment will be the latest available for the enrolling device until this version (excluded). Set for example `28` to automatically required the latest iOS version until (but not including) iOS 28.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
		"macos_min_version": schema.StringAttribute{
			Description:         "The fixed minimum macOS version required for a successful enrollment.",
			MarkdownDescription: "The fixed minimum macOS version required for a successful enrollment.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
		"auto_macos_min_version_until": schema.StringAttribute{
			Description:         "If set, the minimum macOS version required for a successful enrollment will be the latest available for the enrolling device until this version (excluded). Set for example 28 to automatically required the latest macOS version until (but not including) macOS 28.",
			MarkdownDescription: "If set, the minimum macOS version required for a successful enrollment will be the latest available for the enrolling device until this version (excluded). Set for example `28` to automatically required the latest macOS version until (but not including) macOS 28.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
	},
	Optional: true,
	Computed: true,
	Default:  objectdefault.StaticValue(mdeOSVersionEnforcementDefault()),
}

func (r *MDMDEPEnrollmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages MDM ADE/DEP enrollments.",
		MarkdownDescription: "The resource `zentral_mdm_dep_enrollment` manages MDM ADE/DEP enrollments.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description:         "ID of the MDM DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM DEP enrollment.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Description:         "Name of the MDM DEP enrollment as displayed on the device.",
				MarkdownDescription: "Name of the MDM DEP enrollment as displayed on the device.",
				Required:            true,
			},

			"push_certificate_id": schema.Int64Attribute{
				Description:         "ID of the MDM push certificate used by the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM push certificate used by the DEP enrollment.",
				Required:            true,
			},
			"blueprint_id": schema.Int64Attribute{
				Description:         "ID of the MDM blueprint linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM blueprint linked to the DEP enrollment.",
				Optional:            true,
			},

			"acme_issuer_id": schema.StringAttribute{
				Description:         "ID of the optional MDM ACME issuer linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the optional MDM ACME issuer linked to the ODEPenrollment.",
				Optional:            true,
			},
			"scep_issuer_id": schema.StringAttribute{
				Description:         "ID of the MDM SCEP issuer linked to the DEP enrollment.",
				MarkdownDescription: "`ID` of the MDM SCEP issuer linked to the DEP enrollment.",
				Required:            true,
			},

			"enrollment":             enrollmentSecretSchema,
			"authentication":         mdeAuthenticationSchema,
			"extra_admin":            mdeExtraAdminSchema,
			"profile":                mdeProfileSchema,
			"os_version_enforcement": mdeOSVersionEnforcementSchema,
		},
	}
}

func (r *MDMDEPEnrollmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MDMDEPEnrollmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data mdmDEPEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlEnrollment, _, err := r.client.MDMDEPEnrollments.Create(ctx, mdmDEPEnrollmentRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create MDM DEP enrollment, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a MDM DEP enrollment")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentForState(ztlEnrollment))...)
}

func (r *MDMDEPEnrollmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data mdmDEPEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlEnrollment, _, err := r.client.MDMDEPEnrollments.GetByID(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client error",
			fmt.Sprintf("Unable to read MDM DEP enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "read a MDM DEP enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentForState(ztlEnrollment))...)
}

func (r *MDMDEPEnrollmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data mdmDEPEnrollment

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ztlEnrollment, _, err := r.client.MDMDEPEnrollments.Update(ctx, int(data.ID.ValueInt64()), mdmDEPEnrollmentRequestWithState(data))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update MDM DEP enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "updated a MDM DEP enrollment")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, mdmDEPEnrollmentForState(ztlEnrollment))...)
}

func (r *MDMDEPEnrollmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data mdmDEPEnrollment

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.MDMDEPEnrollments.Delete(ctx, int(data.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete MDM DEP enrollment %d, got error: %s", data.ID.ValueInt64(), err),
		)
		return
	}

	tflog.Trace(ctx, "deleted a MDM DEP enrollment")
}

func (r *MDMDEPEnrollmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resourceImportStatePassthroughZentralID(ctx, "MDM DEP enrollment", req, resp)
}
