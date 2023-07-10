// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpenIdServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenIdServerCreate,
		ReadContext:   resourceOpenIdServerRead,
		UpdateContext: resourceOpenIdServerUpdate,
		DeleteContext: resourceOpenIdServerDelete,
		Schema: map[string]*schema.Schema{
			"pid": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tokenIntrospectionClientId": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"tokenIntrospectionClientSecret": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"nodeId": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "Master",
			},
			"moduleId": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "smart_auth",
			},
			"validationJwkText": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"validationJwkFile": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federationRegistrationId": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federationRequestScopes": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "openid profile",
			},
			"federationAuthorizationUrl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"federationTokenUrl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"federationUserInfoUrl": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federationJwkSetUrl": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federationAuthScriptText": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federationAuthScriptFile": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"federationUserMappingScriptText": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"fhirEndpointUrl": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"authWellKnownConfigUrl": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"notes": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"customTokenParams": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"responseType": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"organizationId": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"audience": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
			},
			"archivedAt": &schema.Schema{
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

func resourceOpenIdServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceOpenIdServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceOpenIdServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}

func resourceOpenIdServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags
}
