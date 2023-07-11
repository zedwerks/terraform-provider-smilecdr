// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zed-werks/terraform-smilecdr/client"
)

func SmileCdrProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"baseUrl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http://localhost:9000",
				DefaultFunc: schema.EnvDefaultFunc("SMILECDR_BASE_URL", nil),
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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"smilecdr_openid_client": dataSourceOpenIdClients(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var baseUrl *string

	uVal, ok = d.GetOk("baseUrl")

	if ok {
		base := uVal.(string)
		if !strings.HasSuffix(base, "/") {
			base += "/"
		}
		baseUrl = &base
	}

	if (username != "") && (password != "") {
		c, err := client.NewClient(baseUrl, &username, &password)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Smile CDR client",
				Detail:   "Unable to authenticate user for authenticated Smile CDR client",
			})

			return nil, diags
		}

		return c, diags
	}

	c, err := client.NewClient(baseUrl, nil, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Smile CDR client",
			Detail:   "Unable to create  Smile CDR  client",
		})
		return nil, diags
	}

	return nil, diags
}
