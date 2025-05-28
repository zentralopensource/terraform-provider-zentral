package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

const (
	tfProbeActionHTTPPostBackend             string = "HTTP_POST"
	tfProbeActionSlackIncomingWebhookBackend        = "SLACK_INCOMING_WEBHOOK"
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
	"url":      types.StringType,
	"username": types.StringType,
	"password": types.StringType,
	"headers":  types.SetType{ElemType: types.ObjectType{AttrTypes: probeActionHTTPPostHeaderAttrTypes}},
}

var probeActionSlackIncomingWebhookAttrTypes = map[string]attr.Type{
	"url": types.StringType,
}

func probeActionForState(pa *goztl.ProbeAction) probeAction {
	var hp types.Object
	if pa.HTTPPost != nil {
		// optional username
		var hpu types.String
		if pa.HTTPPost.Username != nil {
			hpu = types.StringValue(*pa.HTTPPost.Username)
		} else {
			hpu = types.StringNull()
		}
		// optional password
		var hpp types.String
		if pa.HTTPPost.Password != nil {
			hpp = types.StringValue(*pa.HTTPPost.Password)
		} else {
			hpp = types.StringNull()
		}
		// optional headers
		hphs := make([]attr.Value, 0)
		for _, h := range pa.HTTPPost.Headers {
			hphs = append(hphs, types.ObjectValueMust(
				probeActionHTTPPostHeaderAttrTypes,
				map[string]attr.Value{
					"name":  types.StringValue(h.Name),
					"value": types.StringValue(h.Value),
				},
			))
		}
		hp = types.ObjectValueMust(
			probeActionHTTPPostAttrTypes,
			map[string]attr.Value{
				"url":      types.StringValue(pa.HTTPPost.URL),
				"username": hpu,
				"password": hpp,
				"headers":  types.SetValueMust(types.ObjectType{AttrTypes: probeActionHTTPPostHeaderAttrTypes}, hphs),
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
			hpBackend := &goztl.ProbeActionHTTPPost{
				URL: hpMap["url"].(types.String).ValueString(),
			}
			// optional username
			hpMapU := hpMap["username"].(types.String)
			if !hpMapU.IsNull() {
				hpBackend.Username = goztl.String(hpMapU.ValueString())
			}
			// optional password
			hpMapP := hpMap["password"].(types.String)
			if !hpMapP.IsNull() {
				hpBackend.Password = goztl.String(hpMapP.ValueString())
			}
			// optional headers
			headers := make([]goztl.ProbeActionHTTPPostHeader, 0)
			for _, hpMapH := range hpMap["headers"].(types.Set).Elements() {
				hpMapHAttrs := hpMapH.(types.Object).Attributes()
				headers = append(headers, goztl.ProbeActionHTTPPostHeader{
					Name:  hpMapHAttrs["name"].(types.String).ValueString(),
					Value: hpMapHAttrs["value"].(types.String).ValueString(),
				})
			}
			hpBackend.Headers = headers
			req.HTTPPost = hpBackend
		}
	}

	if !data.SlackIncomingWebhook.IsNull() {
		siwMap := data.SlackIncomingWebhook.Attributes()
		if siwMap != nil {
			siwBackend := &goztl.ProbeActionSlackIncomingWebhook{
				URL: siwMap["url"].(types.String).ValueString(),
			}
			req.SlackIncomingWebhook = siwBackend
		}
	}

	return req
}
