// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "localhost",
				DefaultFunc: schema.EnvDefaultFunc("SMILECDR_HOST", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SMILECDR_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SMILECDR_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"smilecdr_openid_client": resourceOpenIdClient(),
			"smilecdr_openid_server": resourceOpenIdServer(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"smilecdr_openid_client": dataSourceOpenIdClient(),
			"smilecdr_openid_server": dataSourceOpenIdServer(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics

	return nil, diags
}
