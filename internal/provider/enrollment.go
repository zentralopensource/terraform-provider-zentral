package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

var enrollmentSecretAttrTypes = map[string]attr.Type{
	"secret":                types.StringType,
	"meta_business_unit_id": types.Int64Type,
	"tag_ids":               types.SetType{ElemType: types.Int64Type},
	"serial_numbers":        types.SetType{ElemType: types.StringType},
	"udids":                 types.SetType{ElemType: types.StringType},
	"quota":                 types.Int64Type,
}

func enrollmentSecretForState(e goztl.EnrollmentSecret) types.Object {
	return types.ObjectValueMust(
		enrollmentSecretAttrTypes,
		map[string]attr.Value{
			"secret":                types.StringValue(e.Secret),
			"meta_business_unit_id": types.Int64Value(int64(e.MetaBusinessUnitID)),
			"tag_ids":               int64SetForState(e.TagIDs),
			"serial_numbers":        stringSetForState(e.SerialNumbers),
			"udids":                 stringSetForState(e.UDIDs),
			"quota":                 optionalInt64ForState(e.Quota),
		},
	)
}

func enrollmentSecretRequestWithState(data types.Object) goztl.EnrollmentSecretRequest {
	esMap := data.Attributes()
	return goztl.EnrollmentSecretRequest{
		MetaBusinessUnitID: int(esMap["meta_business_unit_id"].(types.Int64).ValueInt64()),
		TagIDs:             intListWithState(esMap["tag_ids"].(types.Set)),
		SerialNumbers:      stringListWithStateSet(esMap["serial_numbers"].(types.Set)),
		UDIDs:              stringListWithStateSet(esMap["udids"].(types.Set)),
		Quota:              optionalIntWithState(esMap["quota"].(types.Int64)),
	}
}

var enrollmentSecretSchema schema.SingleNestedAttribute = schema.SingleNestedAttribute{
	Description:         "Enrollment settings.",
	MarkdownDescription: "Enrollment settings.",
	Attributes: map[string]schema.Attribute{
		"secret": schema.StringAttribute{
			Description:         "Enrollment secret.",
			MarkdownDescription: "Enrollment secret.",
			Computed:            true,
		},
		"meta_business_unit_id": schema.Int64Attribute{
			Description:         "The ID of the meta business unit the machine will be assigned to at enrollment.",
			MarkdownDescription: "The `ID` of the meta business unit the machine will be assigned to at enrollment.",
			Required:            true,
		},
		"tag_ids": schema.SetAttribute{
			Description:         "The IDs of the tags that the machine will get at enrollment.",
			MarkdownDescription: "The `ID`s of the tags that the machine will get at enrollment.",
			ElementType:         types.Int64Type,
			Optional:            true,
			Computed:            true,
		},
		"serial_numbers": schema.SetAttribute{
			Description:         "The serial numbers the enrollment is restricted to.",
			MarkdownDescription: "The serial numbers the enrollment is restricted to.",
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
		},
		"udids": schema.SetAttribute{
			Description:         "The UDIDs the enrollment is restricted to.",
			MarkdownDescription: "The `UDID`s the enrollment is restricted to.",
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
		},
		"quota": schema.Int64Attribute{
			Description:         "The number of times the enrollment can be used.",
			MarkdownDescription: "The number of times the enrollment can be used.",
			Optional:            true,
		},
	},
	Required: true,
}
