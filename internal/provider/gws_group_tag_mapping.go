package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type gwsGroupTagMapping struct {
	ID           types.String `tfsdk:"id"`
	GroupEmail   types.String `tfsdk:"group_email"`
	ConnectionID types.String `tfsdk:"connection_id"`
	TagIDs       types.Set    `tfsdk:"tag_ids"`
}

func gwsGroupTagMappingForState(gtm *goztl.GWSGroupTagMapping) gwsGroupTagMapping {
	return gwsGroupTagMapping{
		ID:           types.StringValue(gtm.ID),
		GroupEmail:   types.StringValue(gtm.GroupEmail),
		ConnectionID: types.StringValue(gtm.ConnectionID),
		TagIDs:       int64SetForState(gtm.TagIDs),
	}
}

func gwsGroupTagMappingRequestWithState(data gwsGroupTagMapping) *goztl.GWSGroupTagMappingRequest {
	return &goztl.GWSGroupTagMappingRequest{
		GroupEmail:   data.GroupEmail.ValueString(),
		ConnectionID: data.ConnectionID.ValueString(),
		TagIDs:       intListWithState(data.TagIDs),
	}
}
