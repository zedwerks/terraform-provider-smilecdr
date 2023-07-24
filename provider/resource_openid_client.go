// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zedwerks/terraform-smilecdr/provider/helper/validations"
	"github.com/zedwerks/terraform-smilecdr/smilecdr"
)

func resourceOpenIdClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenIdClientCreate,
		ReadContext:   resourceOpenIdClientRead,
		UpdateContext: resourceOpenIdClientUpdate,
		DeleteContext: resourceOpenIdClientDelete,
		Schema: map[string]*schema.Schema{
			"pid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Master",
				ForceNew: true,
			},
			"module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "smart_auth",
				ForceNew: true,
			},
			"access_token_validity_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"allowed_grant_types": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:              schema.HashString,
				ValidateDiagFunc: validateGrantTypes(),
			},
			"auto_approve_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"auto_grant_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"client_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.IsValidClientID,
			},
			"client_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: customValidationFunc,
			},
			"client_secrets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"activation": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
						"expiration": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
					},
				},
			},
			"fixed_scope": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"refresh_token_validity_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  86400,
			},
			"registered_redirect_uris": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"secret_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"secret_client_can_change": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"can_introspect_any_tokens": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_introspect_own_tokens": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"always_require_approval": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_reissue_tokens": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"permissions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validations.IsUserPermission,
						},
						"argument": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"remember_approved_scopes": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"attestation_accepted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"public_jwks_uri": {
				Type:         schema.TypeString,
				Optional:     false,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"archived_at": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"created_by_app_sphere": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDataToOpenIdClient(d *schema.ResourceData) (*smilecdr.OpenIdClient, error) {

	secrets := d.Get("client_secrets").(*schema.Set).List()
	clientSecrets := make([]smilecdr.ClientSecret, len(secrets))

	perms := d.Get("permissions").(*schema.Set).List()
	permissions := make([]smilecdr.UserPermission, len(perms))

	allowedGrantTypes := make([]string, 0)
	allowedGrantTypesData, allowedGrantTypesOk := d.GetOk("allowed_grant_types")

	if allowedGrantTypesOk {
		for _, grant := range allowedGrantTypesData.(*schema.Set).List() {
			allowedGrantTypes = append(allowedGrantTypes, grant.(string))
		}
	}

	autoApproveScopes := make([]string, 0)
	autoApproveScopesData, autoApproveScopesOk := d.GetOk("auto_approve_scopes")

	if autoApproveScopesOk {
		for _, scope := range autoApproveScopesData.(*schema.Set).List() {
			autoApproveScopes = append(autoApproveScopes, scope.(string))
		}
	}
	autoGrantScopes := make([]string, 0)
	autoGrantScopesData, autoGrantScopesOk := d.GetOk("auto_grant_scopes")

	if autoGrantScopesOk {
		for _, scope := range autoGrantScopesData.(*schema.Set).List() {
			autoGrantScopes = append(autoGrantScopes, scope.(string))
		}
	}

	registeredRedirectUris := make([]string, 0)
	registeredRedirectUrisData, registeredRedirectUrisOk := d.GetOk("registered_redirect_uris")

	if registeredRedirectUrisOk {
		for _, uri := range registeredRedirectUrisData.(*schema.Set).List() {
			registeredRedirectUris = append(registeredRedirectUris, uri.(string))
		}
	}

	scopes := make([]string, 0)
	scopesData, scopesOk := d.GetOk("scopes")

	if scopesOk {
		for _, scope := range scopesData.(*schema.Set).List() {
			scopes = append(scopes, scope.(string))
		}
	}

	openidClient := &smilecdr.OpenIdClient{
		Pid:                         d.Get("pid").(int),
		ClientId:                    d.Get("client_id").(string),
		ClientName:                  d.Get("client_name").(string),
		NodeId:                      d.Get("node_id").(string),
		ModuleId:                    d.Get("module_id").(string),
		AccessTokenValiditySeconds:  d.Get("access_token_validity_seconds").(int),
		AllowedGrantTypes:           allowedGrantTypes,
		AutoApproveScopes:           autoApproveScopes,
		AutoGrantScopes:             autoGrantScopes,
		ClientSecrets:               clientSecrets,
		FixedScope:                  d.Get("fixed_scope").(bool),
		RefreshTokenValiditySeconds: d.Get("refresh_token_validity_seconds").(int),
		RegisteredRedirectUris:      registeredRedirectUris,
		Scopes:                      scopes,
		SecretRequired:              d.Get("secret_required").(bool),
		SecretClientCanChange:       d.Get("secret_client_can_change").(bool),
		Enabled:                     d.Get("enabled").(bool),
		CanIntrospectAnyTokens:      d.Get("can_introspect_any_tokens").(bool),
		CanIntrospectOwnTokens:      d.Get("can_introspect_own_tokens").(bool),
		AlwaysRequireApproval:       d.Get("always_require_approval").(bool),
		CanReissueTokens:            d.Get("can_reissue_tokens").(bool),
		Permissions:                 permissions,
		AttestationAccepted:         d.Get("attestation_accepted").(bool),
		PublicJwksUri:               d.Get("public_jwks_uri").(string),
		ArchivedAt:                  d.Get("archived_at").(string),
		CreatedByAppSphere:          d.Get("created_by_app_sphere").(bool),
	}

	return openidClient, nil

}

func resourceOpenIdClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	client, mErr := resourceDataToOpenIdClient(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}

	o, err := c.PostOpenIdClient(*client)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(client.ClientId) // the primary resource identifier. must be unique.
	d.Set("pid", o.Pid)      // the pid is needed for Put requests

	return resourceOpenIdClientRead(ctx, d, m)
}

func resourceOpenIdClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	c := m.(*smilecdr.Client)

	client_id := d.Get("client_id").(string)
	nodeId := d.Get("node_id").(string)
	moduleId := d.Get("module_id").(string)

	openIdClient, err := c.GetOpenIdClient(nodeId, moduleId, client_id)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(openIdClient.ClientId)

	d.Set("pid", openIdClient.Pid)
	d.Set("client_name", openIdClient.ClientName)
	d.Set("node_id", openIdClient.NodeId)
	d.Set("module_id", openIdClient.ModuleId)
	d.Set("access_token_validity_seconds", openIdClient.AccessTokenValiditySeconds)
	d.Set("allowed_grant_types", openIdClient.AllowedGrantTypes)
	d.Set("auto_approve_scopes", openIdClient.AutoApproveScopes)
	d.Set("auto_grant_scopes", openIdClient.AutoGrantScopes)
	d.Set("client_secrets", openIdClient.ClientSecrets)
	d.Set("fixed_scope", openIdClient.FixedScope)
	d.Set("refresh_token_validity_seconds", openIdClient.RefreshTokenValiditySeconds)
	d.Set("registered_redirect_uris", openIdClient.RegisteredRedirectUris)
	d.Set("scopes", openIdClient.Scopes)
	d.Set("secret_required", openIdClient.SecretRequired)
	d.Set("secret_client_can_change", openIdClient.SecretClientCanChange)
	d.Set("enabled", openIdClient.Enabled)
	d.Set("can_introspect_any_tokens", openIdClient.CanIntrospectAnyTokens)
	d.Set("can_introspect_own_tokens", openIdClient.CanIntrospectOwnTokens)
	d.Set("always_require_approval", openIdClient.AlwaysRequireApproval)
	d.Set("can_reissue_tokens", openIdClient.CanReissueTokens)
	d.Set("permissions", openIdClient.Permissions)
	d.Set("attestation_accepted", openIdClient.AttestationAccepted)
	d.Set("public_jwks_uri", openIdClient.PublicJwksUri)
	d.Set("archived_at", openIdClient.ArchivedAt)
	d.Set("created_by_app_sphere", openIdClient.CreatedByAppSphere)
	return diags

}

func resourceOpenIdClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	client, mErr := resourceDataToOpenIdClient(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}

	d.SetId(client.ClientId)

	_, err := c.PutOpenIdClient(*client)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOpenIdClientRead(ctx, d, m)

}

func resourceOpenIdClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	d.Set("archived_at", time.Now().Format(time.RFC3339))

	resourceOpenIdClientUpdate(ctx, d, m)

	d.SetId("")

	return diags
}

// -------------- Validations functions ----------------

func validateClientId() schema.SchemaValidateDiagFunc {
	return validations.IsValidClientID
}

func validateGrantTypes() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		grantTypes := i.([]string)
		return validations.IsValidGrantTypes(grantTypes)
	}
}

func customValidationFunc(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if v.(string) != "expected_value" {
		diag := diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Validation error",
			Detail:   "Expected 'expected_value', but got '" + v.(string) + "'",
		}
		diags = append(diags, diag)
	}

	return diags
}
