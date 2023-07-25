package validations

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ValidateDiagTest(i interface{}, p cty.Path) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Testing Error. Will always throw this.",
			Detail:   "Details!",
		},
	}
}
