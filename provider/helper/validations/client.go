package validations

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var (
	smileCdrGrantTypes = map[string]bool{
		"AUTHORIZATION_CODE": true,
		"IMPLICIT":           true,
		"REFRESH_TOKEN":      true,
		"CLIENT_CREDENTIALS": true,
		"PASSWORD":           true,
		"JWT_BEARER":         true,
	}
)

func IsValidGrantTypes(v interface{}, k cty.Path) diag.Diagnostics {

	var diags diag.Diagnostics

	grants, ok := v.([]string)

	if !ok {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid type for the field",
			Detail:   fmt.Sprintf("Expected a string, but got %T", v),
		}
		diags = append(diags, diag)
	}
	for _, grant := range grants {
		if _, ok := smileCdrGrantTypes[grant]; !ok {
			diag := diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid value for the field",
				Detail:   fmt.Sprintf("Expected valid grant type, but got %s", grant),
			}
			diags = append(diags, diag)
		}
	}
	return diags
}

func IsValidClientID(i interface{}, k cty.Path) diag.Diagnostics {
	value, ok := i.(string)
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid type for the field",
				Detail:   fmt.Sprintf("Expected a string, but got %T", i),
			},
		}
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid type for the field",
				Detail:   fmt.Sprintf("cannot contain whitespace. Got '%s'", value),
			},
		}
	}
	return nil
}
