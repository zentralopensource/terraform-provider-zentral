package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmPushCertificate struct {
	ID              types.Int64  `tfsdk:"id"`
	ProvisioningUID types.String `tfsdk:"provisioning_uid"`
	Name            types.String `tfsdk:"name"`
	Topic           types.String `tfsdk:"topic"`
	Certificate     types.String `tfsdk:"certificate"`
}

func mdmPushCertificateForState(mpc *goztl.MDMPushCertificate) mdmPushCertificate {
	var provisioningUID types.String
	if mpc.ProvisioningUID != nil {
		provisioningUID = types.StringValue(*mpc.ProvisioningUID)
	} else {
		provisioningUID = types.StringNull()
	}

	var topic types.String
	if mpc.Topic != nil {
		topic = types.StringValue(*mpc.Topic)
	} else {
		topic = types.StringNull()
	}

	var certificate types.String
	if mpc.Certificate != nil {
		certificate = types.StringValue(*mpc.Certificate)
	} else {
		certificate = types.StringNull()
	}

	return mdmPushCertificate{
		ID:              types.Int64Value(int64(mpc.ID)),
		ProvisioningUID: provisioningUID,
		Name:            types.StringValue(mpc.Name),
		Topic:           topic,
		Certificate:     certificate,
	}
}
