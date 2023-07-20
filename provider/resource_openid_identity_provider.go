// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zed-werks/terraform-smilecdr/smilecdr"
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
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token_introspection_client_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"token_introspection_client_secret": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "Master",
			},
			"module_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "smart_auth",
			},
			"validation_jwk_text": {
				Type:     schema.TypeString,
				Required: false,
			},
			"validation_jwk_file": {
				Type:     schema.TypeString,
				Required: false,
			},
			"federation_registration_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"federation_request_scopes": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				Default:  "openid profile",
			},
			"federation_authorization_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"federation_token_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"federation_user_info_url": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"federation_jwk_set_url": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"federation_auth_script_text": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"federation_auth_script_file": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"federation_user_mapping_script_text": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"fhir_endpoint_url": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"auth_well_known_config_url": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"custom_token_params": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"response_type": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"audience": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"archived_at": {
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resource2OpenIdIdentityProvider(d *schema.ResourceData) (*smilecdr.OpenIdIdentityProvider, error) {

	openidIdp := &smilecdr.OpenIdIdentityProvider{
		Pid:                             d.Get("pid").(int),
		Name:                            d.Get("name").(string),
		Issuer:                          d.Get("issuer").(string),
		TokenIntrospectionClientId:      d.Get("token_introspection_client_id").(string),
		TokenIntrospectionClientSecret:  d.Get("token_introspection_client_secret").(string),
		NodeId:                          d.Get("node_id").(string),
		ModuleId:                        d.Get("module_id").(string),
		ValidationJwkText:               d.Get("validation_jwk_text").(string),
		ValidationJwkFile:               d.Get("validation_jwk_file").(string),
		FederationRegistrationId:        d.Get("federation_registration_id").(string),
		FederationRequestScopes:         d.Get("federation_request_scopes").(string),
		FederationAuthorizationUrl:      d.Get("federation_authorization_url").(string),
		FederationTokenUrl:              d.Get("federation_token_url").(string),
		FederationUserInfoUrl:           d.Get("federation_user_info_url").(string),
		FederationJwkSetUrl:             d.Get("federation_jwk_set_url").(string),
		FederationAuthScriptText:        d.Get("federation_auth_script_text").(string),
		FederationAuthScriptFile:        d.Get("federation_auth_script_file").(string),
		FederationUserMappingScriptText: d.Get("federation_user_mapping_script_text").(string),
		FhirEndpointUrl:                 d.Get("fhir_endpoint_url").(string),
		AuthWellKnownConfigUrl:          d.Get("auth_well_known_config_url").(string),
		Notes:                           d.Get("notes").(string),
		CustomTokenParams:               d.Get("custom_token_params").(string),
		ResponseType:                    d.Get("response_type").(string),
		OrganizationId:                  d.Get("organization_id").(string),
		Audience:                        d.Get("audience").(string),
		ArchivedAt:                      d.Get("archived_at").(string),
	}

	return openidIdp, nil
}

func resourceOpenIdIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	idp, mErr := resource2OpenIdIdentityProvider(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}

	o, err := c.PostOpenIdIdentityProvider(*idp)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(idp.Issuer) // the primary resource identifier. must be unique.
	d.Set("pid", o.Pid) // the pid is needed for Put requests

	return resourceOpenIdIdentityProviderRead(ctx, d, m)
}

func resourceOpenIdIdentityProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	c := m.(*smilecdr.Client)

	pid := d.Get("pid").(string)
	nodeId := d.Get("node_id").(string)
	moduleId := d.Get("module_id").(string)

	provider, err := c.GetOpenIdIdentityProvider(nodeId, moduleId, pid)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(provider.Issuer)

	d.Set("pid", provider.Pid)
	d.Set("name", provider.Name)
	d.Set("issuer", provider.Issuer)
	d.Set("token_introspection_client_id", provider.TokenIntrospectionClientId)
	d.Set("token_introspection_client_secret", provider.TokenIntrospectionClientSecret)
	d.Set("node_id", provider.NodeId)
	d.Set("module_id", provider.ModuleId)
	d.Set("validation_jwk_text", provider.ValidationJwkText)
	d.Set("validation_jwk_file", provider.ValidationJwkFile)
	d.Set("federation_registration_id", provider.FederationRegistrationId)
	d.Set("federation_request_scopes", provider.FederationRequestScopes)
	d.Set("federation_authorization_url", provider.FederationAuthorizationUrl)
	d.Set("federation_token_url", provider.FederationTokenUrl)
	d.Set("federation_user_info_url", provider.FederationUserInfoUrl)
	d.Set("federation_jwk_set_url", provider.FederationJwkSetUrl)
	d.Set("federation_auth_script_text", provider.FederationAuthScriptText)
	d.Set("federation_auth_script_file", provider.FederationAuthScriptFile)
	d.Set("federation_user_mapping_script_text", provider.FederationUserMappingScriptText)
	d.Set("fhir_endpoint_url", provider.FhirEndpointUrl)
	d.Set("auth_well_known_config_url", provider.AuthWellKnownConfigUrl)
	d.Set("notes", provider.Notes)
	d.Set("custom_token_params", provider.CustomTokenParams)
	d.Set("response_type", provider.ResponseType)
	d.Set("organization_id", provider.OrganizationId)
	d.Set("audience", provider.Audience)
	d.Set("archived_at", provider.ArchivedAt)

	return diags
}

func resourceOpenIdIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	provider, mErr := resource2OpenIdIdentityProvider(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}
	d.SetId(provider.Name)

	_, err := c.PutOpenIdIdentityProvider(*provider)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOpenIdIdentityProviderRead(ctx, d, m)

}

func resourceOpenIdIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	d.Set("archived_at", time.Now().Format(time.RFC3339))

	resourceOpenIdIdentityProviderUpdate(ctx, d, m)

	d.SetId("")

	return diags
}
