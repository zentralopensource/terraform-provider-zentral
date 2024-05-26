package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	tagIDs := make([]attr.Value, 0)
	for _, tagID := range mmep.TagIDs {
		tagIDs = append(tagIDs, types.Int64Value(int64(tagID)))
	}

	return monolithManifestEnrollmentPackage{
		ID:           types.Int64Value(int64(mmep.ID)),
		ManifestID:   types.Int64Value(int64(mmep.ManifestID)),
		Builder:      types.StringValue(mmep.Builder),
		EnrollmentID: types.Int64Value(int64(mmep.EnrollmentID)),
		Version:      types.Int64Value(int64(mmep.Version)),
		TagIDs:       types.SetValueMust(types.Int64Type, tagIDs),
	}
}

func monolithManifestEnrollmentPackageRequestWithState(data monolithManifestEnrollmentPackage) *goztl.MonolithManifestEnrollmentPackageRequest {
	tagIDs := make([]int, 0)
	for _, tagID := range data.TagIDs.Elements() { // nil if null or unknown → no iterations
		tagIDs = append(tagIDs, int(tagID.(types.Int64).ValueInt64()))
	}

	return &goztl.MonolithManifestEnrollmentPackageRequest{
		ManifestID:   int(data.ManifestID.ValueInt64()),
		Builder:      data.Builder.ValueString(),
		EnrollmentID: int(data.EnrollmentID.ValueInt64()),
		TagIDs:       tagIDs,
	}
}
