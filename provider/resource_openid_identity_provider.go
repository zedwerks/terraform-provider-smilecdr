// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpenIdIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenIdIdentityProviderCreate,
		ReadContext:   resourceOpenIdIdentityProviderRead,
		UpdateContext: resourceOpenIdIdentityProviderUpdate,
		DeleteContext: resourceOpenIdIdentityProviderDelete,
		Schema: map[string]*schema.Schema{
			"pid": {
				Type:     schema.TypeInt,
				Required: false,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tokenIntrospectionClientId": {
				Type:     schema.TypeString,
				Required: false,
			},
			"tokenIntrospectionClientSecret": {
				Type:     schema.TypeString,
				Required: false,
			},
			"nodeId": {
				Type:     schema.TypeString,
				Required: false,
				Default:  "Master",
			},
			"moduleId": {
				Type:     schema.TypeString,
				Required: false,
				Default:  "smart_auth",
			},
			"validationJwkText": {
				Type:     schema.TypeString,
				Required: false,
			},
			"validationJwkFile": {
				Type:     schema.TypeString,
				Required: false,
			},
			"federationRegistrationId": {
				Type:     schema.TypeString,
				Required: false,
			},
			"federationRequestScopes": {
				Type:     schema.TypeString,
				Required: false,
				Default:  "openid profile",
			},
			"federationAuthorizationUrl": {
				Type:     schema.TypeString,
				Required: true,
			},
			"federationTokenUrl": {
				Type:     schema.TypeString,
				Required: true,
			},
			"federationUserInfoUrl": {
				Type:     schema.TypeString,
				Required: false,
			},
			"federationJwkSetUrl": {
				Type:     schema.TypeString,
				Required: false,
			},
			"federationAuthScriptText": {
				Type:     schema.TypeString,
				Required: false,
			},
			"federationAuthScriptFile": {
				Type:     schema.TypeString,
				Required: false,
			},
			"federationUserMappingScriptText": {
				Type:     schema.TypeString,
				Required: false,
			},
			"fhirEndpointUrl": {
				Type:     schema.TypeString,
				Required: false,
			},
			"authWellKnownConfigUrl": {
				Type:     schema.TypeString,
				Required: false,
			},
			"notes": {
				Type:     schema.TypeString,
				Required: false,
			},
			"customTokenParams": {
				Type:     schema.TypeString,
				Required: false,
			},
			"responseType": {
				Type:     schema.TypeString,
				Required: false,
			},
			"organizationId": {
				Type:     schema.TypeString,
				Required: false,
			},
			"audience": {
				Type:     schema.TypeString,
				Required: false,
			},
			"archivedAt": {
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

func resourceOpenIdIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceOpenIdIdentityProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceOpenIdIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceOpenIdIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
