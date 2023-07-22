package utils

import (
	"fmt"
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

func ValidateGrantType(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	// Check if the value is in the set of acceptable values
	if _, ok := smileCdrGrantTypes[value]; !ok {
		errs = append(errs, fmt.Errorf("invalid grant type. Got %s", value))
	}
	return warns, errs
}
