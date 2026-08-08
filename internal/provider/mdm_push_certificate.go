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
	return mdmPushCertificate{
		ID:              types.Int64Value(int64(mpc.ID)),
		ProvisioningUID: optionalStringForState(mpc.ProvisioningUID),
		Name:            types.StringValue(mpc.Name),
		Topic:           optionalStringForState(mpc.Topic),
		Certificate:     optionalStringForState(mpc.Certificate),
	}
}
