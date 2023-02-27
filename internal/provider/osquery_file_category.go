package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	filePaths := make([]attr.Value, 0)
	for _, filePath := range ofc.FilePaths {
		filePaths = append(filePaths, types.StringValue(filePath))
	}

	excludePaths := make([]attr.Value, 0)
	for _, excludePath := range ofc.ExcludePaths {
		excludePaths = append(excludePaths, types.StringValue(excludePath))
	}

	filePathsQueries := make([]attr.Value, 0)
	for _, filePathsQuery := range ofc.FilePathsQueries {
		filePathsQueries = append(filePathsQueries, types.StringValue(filePathsQuery))
	}

	return osqueryFileCategory{
		ID:               types.Int64Value(int64(ofc.ID)),
		Name:             types.StringValue(ofc.Name),
		Description:      types.StringValue(ofc.Description),
		FilePaths:        types.SetValueMust(types.StringType, filePaths),
		ExcludePaths:     types.SetValueMust(types.StringType, excludePaths),
		FilePathsQueries: types.SetValueMust(types.StringType, filePathsQueries),
		AccessMonitoring: types.BoolValue(ofc.AccessMonitoring),
	}
}

func osqueryFileCategoryRequestWithState(data osqueryFileCategory) *goztl.OsqueryFileCategoryRequest {
	filePaths := make([]string, 0)
	for _, filePath := range data.FilePaths.Elements() { // nil if null or unknown → no iterations
		filePaths = append(filePaths, filePath.(types.String).ValueString())
	}

	excludePaths := make([]string, 0)
	for _, excludePath := range data.ExcludePaths.Elements() { // nil if null or unknown → no iterations
		excludePaths = append(excludePaths, excludePath.(types.String).ValueString())
	}

	filePathsQueries := make([]string, 0)
	for _, filePathsQuery := range data.FilePathsQueries.Elements() { // nil if null or unknown → no iterations
		filePathsQueries = append(filePathsQueries, filePathsQuery.(types.String).ValueString())
	}

	return &goztl.OsqueryFileCategoryRequest{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		FilePaths:        filePaths,
		ExcludePaths:     excludePaths,
		FilePathsQueries: filePathsQueries,
		AccessMonitoring: data.AccessMonitoring.ValueBool(),
	}
}
