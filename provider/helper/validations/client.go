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
