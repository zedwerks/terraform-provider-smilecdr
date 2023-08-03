// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

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
		c := smilecdr.NewClient(baseUrl, username, password)

		return c, diags
	}

	return nil, diags
}
