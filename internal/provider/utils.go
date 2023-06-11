package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceImportStatePassthroughZentralID(ctx context.Context, name string, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ztlID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Zentral %s ID must be an integer", name),
		)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.Int64Value(ztlID))...)
	}
}

var (
	uuidMatchRe = regexp.MustCompile(`[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}`)
)

func resourceImportStatePassthroughZentralUUID(ctx context.Context, name string, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	matched := uuidMatchRe.MatchString(req.ID)
	if !matched {
		resp.Diagnostics.AddError(
			"Invalid resource ID",
			fmt.Sprintf("Zentral %s ID must be a valid UUID", name),
		)
	} else {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	}
}
