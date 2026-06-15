package provider

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zentralopensource/goztl"
)

type mdmPackage struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	SourceURI      types.String `tfsdk:"source_uri"`
	SHA256         types.String `tfsdk:"sha256"`
	Type           types.String `tfsdk:"type"`
	Size           types.Int64  `tfsdk:"size"`
	Filename       types.String `tfsdk:"filename"`
	ProductID      types.String `tfsdk:"product_id"`
	ProductVersion types.String `tfsdk:"product_version"`
	Bundles        types.String `tfsdk:"bundles"`
	Manifest       types.String `tfsdk:"manifest"`
}

// bundles and manifest are opaque JSON blobs computed server-side from the
// uploaded file. They are marshalled into a Computed types.String — operators
// never write them, so plan-stability only requires that the Create and Read
// code paths produce identical bytes for the same server payload.

func mdmPackageForState(mp *goztl.MDMPackage) (mdmPackage, error) {
	bundles, err := json.Marshal(mp.Bundles)
	if err != nil {
		return mdmPackage{}, err
	}
	manifest, err := json.Marshal(mp.Manifest)
	if err != nil {
		return mdmPackage{}, err
	}
	return mdmPackage{
		ID:             types.StringValue(mp.ID),
		Name:           types.StringValue(mp.Name),
		Description:    types.StringValue(mp.Description),
		SourceURI:      types.StringValue(mp.SourceURI),
		SHA256:         types.StringValue(mp.SHA256),
		Type:           types.StringValue(mp.Type),
		Size:           types.Int64Value(mp.Size),
		Filename:       types.StringValue(mp.Filename),
		ProductID:      types.StringValue(mp.ProductID),
		ProductVersion: types.StringValue(mp.ProductVersion),
		Bundles:        types.StringValue(string(bundles)),
		Manifest:       types.StringValue(string(manifest)),
	}, nil
}

// Create takes name + description + source_uri + sha256. The server downloads
// the file at source_uri, verifies the hash, and derives everything else.
func mdmPackageCreateRequestWithState(data mdmPackage) *goztl.MDMPackageCreateRequest {
	return &goztl.MDMPackageCreateRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		SourceURI:   data.SourceURI.ValueString(),
		SHA256:      data.SHA256.ValueString(),
	}
}

// Update only accepts name + description. source_uri and sha256 are immutable
// post-create and gated at the schema level by RequiresReplace; this builder
// intentionally omits them.
func mdmPackageUpdateRequestWithState(data mdmPackage) *goztl.MDMPackageUpdateRequest {
	return &goztl.MDMPackageUpdateRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}
}
