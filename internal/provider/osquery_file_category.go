package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type osqueryFileCategory struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	FilePaths        types.Set    `tfsdk:"file_paths"`
	ExcludePaths     types.Set    `tfsdk:"exclude_paths"`
	FilePathsQueries types.Set    `tfsdk:"file_paths_queries"`
	AccessMonitoring types.Bool   `tfsdk:"access_monitoring"`
}

func osqueryFileCategoryForState(ofc *goztl.OsqueryFileCategory) osqueryFileCategory {
	return osqueryFileCategory{
		ID:               types.Int64Value(int64(ofc.ID)),
		Name:             types.StringValue(ofc.Name),
		Description:      types.StringValue(ofc.Description),
		FilePaths:        stringSetForState(ofc.FilePaths),
		ExcludePaths:     stringSetForState(ofc.ExcludePaths),
		FilePathsQueries: stringSetForState(ofc.FilePathsQueries),
		AccessMonitoring: types.BoolValue(ofc.AccessMonitoring),
	}
}

func osqueryFileCategoryRequestWithState(data osqueryFileCategory) *goztl.OsqueryFileCategoryRequest {
	return &goztl.OsqueryFileCategoryRequest{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		FilePaths:        stringListWithStateSet(data.FilePaths),
		ExcludePaths:     stringListWithStateSet(data.ExcludePaths),
		FilePathsQueries: stringListWithStateSet(data.FilePathsQueries),
		AccessMonitoring: data.AccessMonitoring.ValueBool(),
	}
}
