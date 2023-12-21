// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zedwerks/terraform-smilecdr/provider/helper/validations"
	"github.com/zedwerks/terraform-smilecdr/smilecdr"
)

func resourceOpenIdIdentityProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenIdIdentityProviderCreate,
		ReadContext:   resourceOpenIdIdentityProviderRead,
		UpdateContext: resourceOpenIdIdentityProviderUpdate,
		DeleteContext: resourceOpenIdIdentityProviderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOpenIdIdentityProviderImport,
		},
		Schema: map[string]*schema.Schema{
			"pid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "An OpenID Identity Provider",
			},
			"issuer": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
			},
			"token_introspection_client_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"token_introspection_client_secret": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.StringLenBetween(9, 512)),
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Master",
			},
			"module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "smart_auth",
			},
			"validation_jwk_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"validation_jwk_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"federation_registration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"federation_request_scopes": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "openid profile",
			},
			"federation_authorization_url": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
			},
			"federation_token_url": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
			},
			"federation_user_info_url": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
			},
			"federation_jwk_set_url": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
			},
			"federation_auth_script_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"federation_user_mapping_script_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"archived_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
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
		FederationUserMappingScriptText: d.Get("federation_user_mapping_script_text").(string),
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

	o, err := c.PostOpenIdIdentityProvider(ctx, *idp)

	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error creating identity provider",
		})
		return diags
	}

	d.Set("federation_registration_id", o.FederationRegistrationId) // set the computed value
	d.SetId(o.Issuer)                                               // the primary resource identifier. must be unique.
	d.Set("pid", o.Pid)                                             // the pid is needed for Put requests

	return resourceOpenIdIdentityProviderRead(ctx, d, m)
}

func resourceOpenIdIdentityProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	c := m.(*smilecdr.Client)

	nodeId := d.Get("node_id").(string)
	moduleId := d.Get("module_id").(string)
	issuerUrl := d.Get("issuer").(string)

	provider, err := c.GetOpenIdIdentityProvider(ctx, nodeId, moduleId, issuerUrl)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(provider.Issuer) // Issuer must be unique in the system

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
	d.Set("federation_user_mapping_script_text", provider.FederationUserMappingScriptText)
	d.Set("archived_at", provider.ArchivedAt)

	return diags
}

func resourceOpenIdIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	provider, cErr := resource2OpenIdIdentityProvider(d)
	if cErr != nil {
		diags := diag.FromErr(cErr)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error converting identity provider resource information to data model",
		})
		return diags
	}
	d.SetId(provider.Issuer)

	_, err := c.PutOpenIdIdentityProvider(ctx, *provider)

	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error updating identity provider",
		})
		return diags
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

func resourceOpenIdIdentityProviderImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(*smilecdr.Client)

	parts := strings.Split(d.Id(), "?issuer_url=")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import. supported import formats: {{nodeId}}/{{moduleId}}?issuer_url={{issuerUrl}}")
	}

	moduleParts := strings.Split(parts[0], "/")
	if len(moduleParts) != 2 {
		return nil, fmt.Errorf("invalid import. supported import formats: {{nodeId}}/{{moduleId}}?issuer_url={{issuerUrl}}")
	}

	_, err := c.GetOpenIdIdentityProvider(ctx, moduleParts[0], moduleParts[1], parts[2])
	if err != nil {
		return nil, err
	}

	d.Set("node_id", moduleParts[0])
	d.Set("module_id", moduleParts[1])
	d.Set("issuer", parts[2])

	diagnostics := resourceOpenIdIdentityProviderRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
