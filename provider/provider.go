// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zedwerks/terraform-smilecdr/smilecdr"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "http://localhost:9000",
				DefaultFunc:  schema.EnvDefaultFunc("SMILECDR_BASE_URL", nil),
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SMILECDR_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SMILECDR_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"smilecdr_openid_client":            resourceOpenIdClient(),
			"smilecdr_openid_identity_provider": resourceOpenIdIdentityProvider(),
			"smilecdr_smart_outbound_security":  resourceSmartOutboundSecurity(),
			"smilecdr_smart_inbound_security":   resourceSmartInboundSecurity(),
			"smilecdr_module_config":            resourceModuleConfig(),
			"smilecdr_user":                     resourceUser(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics
	var baseUrl string

	username := d.Get("username").(string)
	password := d.Get("password").(string)
	baseUrl = d.Get("base_url").(string)

	if (baseUrl != "") && (username != "") && (password != "") {
		c := smilecdr.NewClient(ctx, baseUrl, username, password)

		return c, diags
	}

	return nil, diags
}

func suppressSensitiveDataDiff(k, old, new string, d *schema.ResourceData) (bool, error) {
	// Your custom diff logic here
	log.Printf("Diff function called for key: %s, old: %s, new: %s\n", k, old, new)

	// For example, let's suppress the diff for the 'password' attribute
	if k == "password" {
		log.Println("Suppressing diff for sensitive attribute 'password'")
		return true, nil
	}
	if k == "secret" {
		log.Println("Suppressing diff for sensitive attribute 'secret'")
		return true, nil
	}

	// Continue with the default behavior for other attributes
	return false, nil
}
