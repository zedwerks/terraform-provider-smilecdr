package utils

import (
	"fmt"
	"net/url"
	"regexp"
)

func ValidateClientId(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func ValidateUrl(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected name to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("name cannot contain whitespace. Got %s", value))
		return warns, errs
	}

	parsedUrl, err := url.Parse(value)
	if err != nil {
		errs = append(errs, fmt.Errorf("url is not valid. Got %s", value))
		return warns, errs
	}
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		errs = append(errs, fmt.Errorf("url scheme must be http or https. Got %s", parsedUrl.Scheme))
		return warns, errs
	}
	if parsedUrl.Host == "" {
		errs = append(errs, fmt.Errorf("url must contain a host. Got %s", parsedUrl.Host))
		return warns, errs
	}

	return warns, errs
}
