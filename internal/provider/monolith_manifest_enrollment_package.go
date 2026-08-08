package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type monolithManifestEnrollmentPackage struct {
	ID           types.Int64  `tfsdk:"id"`
	ManifestID   types.Int64  `tfsdk:"manifest_id"`
	Builder      types.String `tfsdk:"builder"`
	EnrollmentID types.Int64  `tfsdk:"enrollment_id"`
	Version      types.Int64  `tfsdk:"version"`
	TagIDs       types.Set    `tfsdk:"tag_ids"`
}

func monolithManifestEnrollmentPackageForState(mmep *goztl.MonolithManifestEnrollmentPackage) monolithManifestEnrollmentPackage {
	return monolithManifestEnrollmentPackage{
		ID:           types.Int64Value(int64(mmep.ID)),
		ManifestID:   types.Int64Value(int64(mmep.ManifestID)),
		Builder:      types.StringValue(mmep.Builder),
		EnrollmentID: types.Int64Value(int64(mmep.EnrollmentID)),
		Version:      types.Int64Value(int64(mmep.Version)),
		TagIDs:       int64SetForState(mmep.TagIDs),
	}
}

func monolithManifestEnrollmentPackageRequestWithState(data monolithManifestEnrollmentPackage) *goztl.MonolithManifestEnrollmentPackageRequest {
	return &goztl.MonolithManifestEnrollmentPackageRequest{
		ManifestID:   int(data.ManifestID.ValueInt64()),
		Builder:      data.Builder.ValueString(),
		EnrollmentID: int(data.EnrollmentID.ValueInt64()),
		TagIDs:       intListWithState(data.TagIDs),
	}
}
