package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfProbeActionHTTPPostBackend             string = "HTTP_POST"
	tfProbeActionSlackIncomingWebhookBackend string = "SLACK_INCOMING_WEBHOOK"
)

type probeAction struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Backend              types.String `tfsdk:"backend"`
	HTTPPost             types.Object `tfsdk:"http_post"`
	SlackIncomingWebhook types.Object `tfsdk:"slack_incoming_webhook"`
}

var probeActionHTTPPostHeaderAttrTypes = map[string]attr.Type{
	"name":  types.StringType,
	"value": types.StringType,
}

var probeActionHTTPPostAttrTypes = map[string]attr.Type{
	"url":                types.StringType,
	"username":           types.StringType,
	"password":           types.StringType,
	"headers":            types.SetType{ElemType: types.ObjectType{AttrTypes: probeActionHTTPPostHeaderAttrTypes}},
	"cel_transformation": types.StringType,
}

var probeActionSlackIncomingWebhookAttrTypes = map[string]attr.Type{
	"url": types.StringType,
}

func probeActionForState(pa *goztl.ProbeAction) probeAction {
	var hp types.Object
	if pa.HTTPPost != nil {
		hp = types.ObjectValueMust(
			probeActionHTTPPostAttrTypes,
			map[string]attr.Value{
				"url":                types.StringValue(pa.HTTPPost.URL),
				"username":           optionalStringForState(pa.HTTPPost.Username),
				"password":           optionalStringForState(pa.HTTPPost.Password),
				"headers":            headersForState(pa.HTTPPost.Headers),
				"cel_transformation": optionalStringForState(pa.HTTPPost.CELTransformation),
			},
		)
	} else {
		hp = types.ObjectNull(probeActionHTTPPostAttrTypes)
	}

	var siw types.Object
	if pa.SlackIncomingWebhook != nil {
		siw = types.ObjectValueMust(
			probeActionSlackIncomingWebhookAttrTypes,
			map[string]attr.Value{
				"url": types.StringValue(pa.SlackIncomingWebhook.URL),
			},
		)
	} else {
		siw = types.ObjectNull(probeActionSlackIncomingWebhookAttrTypes)
	}

	return probeAction{
		ID:                   types.StringValue(pa.ID),
		Name:                 types.StringValue(pa.Name),
		Description:          types.StringValue(pa.Description),
		Backend:              types.StringValue(pa.Backend),
		HTTPPost:             hp,
		SlackIncomingWebhook: siw,
	}
}

func probeActionRequestWithState(data probeAction) *goztl.ProbeActionRequest {
	req := &goztl.ProbeActionRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Backend:     data.Backend.ValueString(),
	}

	if !data.HTTPPost.IsNull() {
		hpMap := data.HTTPPost.Attributes()
		if hpMap != nil {
			req.HTTPPost = &goztl.ProbeActionHTTPPost{
				URL:               hpMap["url"].(types.String).ValueString(),
				Username:          optionalStringWithState(hpMap["username"].(types.String)),
				Password:          optionalStringWithState(hpMap["password"].(types.String)),
				Headers:           headersWithState(hpMap["headers"].(types.Set)),
				CELTransformation: optionalStringWithState(hpMap["cel_transformation"].(types.String)),
			}
		}
	}

	if !data.SlackIncomingWebhook.IsNull() {
		siwMap := data.SlackIncomingWebhook.Attributes()
		if siwMap != nil {
			req.SlackIncomingWebhook = &goztl.ProbeActionSlackIncomingWebhook{
				URL: siwMap["url"].(types.String).ValueString(),
			}
		}
	}

	return req
}
