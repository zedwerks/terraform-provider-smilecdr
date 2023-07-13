// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpenIdIdentityProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenIdIdentityProviderRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"token_introspection_client_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"tokenIntrospectionClientSecret": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"node_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "Master",
			},
			"module_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "smart_auth",
			},
			"validation_jwk_text": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"validation_jwk_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federation_registration_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federation_request_scopes": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "openid profile",
			},
			"federation_authorization_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"federation_token_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"federation_user_info_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federation_jwk_set_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federation_auth_script_text": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federation_auth_script_file": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federation_user_mapping_script_text": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"fhir_endpoint_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"auth_well_known_config_url": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"custom_token_params": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"response_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"audience": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"archived_at": &schema.Schema{
				Type:         schema.TypeString,
				Required:     false,
				ValidateFunc: validation.IsRFC3339Time,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceOpenIdIdentityProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
