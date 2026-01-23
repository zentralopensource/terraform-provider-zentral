package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmEnrollmentCutomView struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	HTML                   types.String `tfsdk:"html"`
	HTMLFile               types.String `tfsdk:"html_file"`
	RequiresAuthentication types.Bool   `tfsdk:"requires_authentication"`
}

type mdmEnrollmentCustomViewDataSourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Description            types.String `tfsdk:"description"`
	HTML                   types.String `tfsdk:"html"`
	RequiresAuthentication types.Bool   `tfsdk:"requires_authentication"`
}

func mdmEnrollmentCustomViewDSForState(v *goztl.MDMEnrollmentCustomView) mdmEnrollmentCustomViewDataSourceModel {
	return mdmEnrollmentCustomViewDataSourceModel{
		ID:                     types.StringValue(v.ID),
		Name:                   types.StringValue(v.Name),
		Description:            types.StringValue(v.Description),
		HTML:                   types.StringValue(v.HTML),
		RequiresAuthentication: types.BoolValue(v.RequiresAuthentication),
	}
}

func mdmEnrollmentCustomViewForState(customView *goztl.MDMEnrollmentCustomView) mdmEnrollmentCutomView {
	return mdmEnrollmentCutomView{
		ID:                     types.StringValue(customView.ID),
		Name:                   types.StringValue(customView.Name),
		Description:            types.StringValue(customView.Description),
		HTML:                   types.StringValue(customView.HTML),
		RequiresAuthentication: types.BoolValue(customView.RequiresAuthentication),
	}
}

func resolveMDMEnrollmentCustomViewHTML(data mdmEnrollmentCutomView) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Prefer direct html when set
	if !data.HTML.IsNull() && data.HTML.ValueString() != "" {
		return data.HTML.ValueString(), diags
	}

	// Fallback to html_file when set
	if !data.HTMLFile.IsNull() && data.HTMLFile.ValueString() != "" {
		b, err := os.ReadFile(data.HTMLFile.ValueString())
		if err != nil {
			diags.AddError(
				"Cannot read html_file",
				fmt.Sprintf("Unable to read file %q: %s", data.HTMLFile.ValueString(), err),
			)
			return "", diags
		}
		return string(b), diags
	}

	diags.AddError(
		"Missing HTML",
		"Either `html` or `html_file` must be set.",
	)
	return "", diags
}

func mdmEnrollmentCustomViewRequestWithState(data mdmEnrollmentCutomView) (*goztl.MDMEnrollmentCustomViewRequest, diag.Diagnostics) {
	html, diags := resolveMDMEnrollmentCustomViewHTML(data)
	if diags.HasError() {
		return nil, diags
	}
	mdmEnrollmentCustomViewRequest := &goztl.MDMEnrollmentCustomViewRequest{
		Name:                   data.Name.ValueString(),
		Description:            data.Description.ValueString(),
		HTML:                   html,
		RequiresAuthentication: data.RequiresAuthentication.ValueBool(),
	}
	return mdmEnrollmentCustomViewRequest, diags
}

func resolveHTMLFromConfig(ctx context.Context, data mdmEnrollmentCutomView) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if !data.HTML.IsNull() && data.HTML.ValueString() != "" {
		return data.HTML.ValueString(), diags
	}

	if !data.HTMLFile.IsNull() && data.HTMLFile.ValueString() != "" {
		b, err := os.ReadFile(data.HTMLFile.ValueString())
		if err != nil {
			diags.AddError(
				"Cannot read html_file",
				fmt.Sprintf("Unable to read file %q: %s", data.HTMLFile.ValueString(), err),
			)
			return "", diags
		}
		return string(b), diags
	}

	// Sollte durch ExactlyOneOf eigentlich nicht passieren, aber defensiv:
	diags.AddError("Missing HTML", "Either `html` or `html_file` must be set.")
	return "", diags
}
