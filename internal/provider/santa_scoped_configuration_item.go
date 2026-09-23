package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The scoped client modes and the scoped path regexes share their base server side: the
// configuration, a name unique in it, a description, and the scope fields of a rule.
func santaScopedConfigurationItemAttributes(item string) map[string]schema.Attribute {
	scopeStrings := func(description string) schema.SetAttribute {
		return schema.SetAttribute{
			Description:         description,
			MarkdownDescription: description,
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
			Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1), trimmedStringValidator),
			},
		}
	}
	scopeTags := func(description, markdownDescription string) schema.SetAttribute {
		return schema.SetAttribute{
			Description:         description,
			MarkdownDescription: markdownDescription,
			ElementType:         types.Int64Type,
			Optional:            true,
			Computed:            true,
			Default:             setdefault.StaticValue(types.SetValueMust(types.Int64Type, []attr.Value{})),
		}
	}

	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description:         fmt.Sprintf("ID of the Santa scoped %s.", item),
			MarkdownDescription: fmt.Sprintf("`ID` of the Santa scoped %s.", item),
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"configuration_id": schema.Int64Attribute{
			Description:         "ID of the Santa configuration. Changing it replaces the entry.",
			MarkdownDescription: "`ID` of the Santa configuration. Changing it replaces the entry.",
			Required:            true,
			PlanModifiers: []planmodifier.Int64{
				// the server refuses to move an entry to another configuration
				int64planmodifier.RequiresReplace(),
			},
		},
		"name": schema.StringAttribute{
			Description:         fmt.Sprintf("Name of the Santa scoped %s, unique in the configuration.", item),
			MarkdownDescription: fmt.Sprintf("Name of the Santa scoped %s, unique in the configuration.", item),
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 256),
				trimmedStringValidator,
			},
		},
		"description": schema.StringAttribute{
			Description:         fmt.Sprintf("Description of the Santa scoped %s.", item),
			MarkdownDescription: fmt.Sprintf("Description of the Santa scoped %s.", item),
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
			Validators: []validator.String{
				trimmedStringValidator,
			},
		},
		"primary_users": scopeStrings(
			"The primary users of the machines in scope.",
		),
		"excluded_primary_users": scopeStrings(
			"The primary users of the machines out of scope.",
		),
		"serial_numbers": scopeStrings(
			"The serial numbers of the machines in scope.",
		),
		"excluded_serial_numbers": scopeStrings(
			"The serial numbers of the machines out of scope.",
		),
		"tag_ids": scopeTags(
			"The IDs of the tags of the machines in scope.",
			"The `ID`s of the tags of the machines in scope.",
		),
		"excluded_tag_ids": scopeTags(
			"The IDs of the tags of the machines out of scope.",
			"The `ID`s of the tags of the machines out of scope.",
		),
	}
}
