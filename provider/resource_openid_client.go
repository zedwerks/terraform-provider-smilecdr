// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"
	"fmt"
	"time"

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
			},
			"module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "smart_auth",
			},
			"client_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.IsValidClientID,
			},
			"client_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"access_token_validity_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"allowed_grant_types": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validations.ValidateDiagFunc(validation.StringInSlice([]string{"AUTHORIZATION_CODE", "IMPLICIT", "REFRESH_TOKEN", "CLIENT_CREDENTIALS", "PASSWORD", "JWT_BEARER"}, false)),
				},
			},
			"always_require_approval": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"attestation_accepted": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"auto_approve_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"auto_grant_scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"can_reissue_tokens": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"client_secrets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validations.ValidateDiagFunc(validation.StringLenBetween(10, 256)),
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"activation": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsRFC3339Time),
						},
						"expiration": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsRFC3339Time),
						},
					},
				},
			},
			"created_by_app_sphere": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"fixed_scope": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"jwks_url": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
			},
			"permissions": {
				Type:     schema.TypeSet,
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
			"public_jwks": {
				Type:     schema.TypeString,
				Optional: true,
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
					Type:             schema.TypeString,
					ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
				},
			},
			"remember_approved_scopes": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"secret_client_can_change": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"secret_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"archived_at": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsRFC3339Time),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func flattenClientSecrets(clientSecrets []smilecdr.ClientSecret) []interface{} {
	secrets := make([]interface{}, len(clientSecrets))

	for i, s := range clientSecrets {
		secrets[i] = map[string]interface{}{
			"secret":      s.Secret,
			"description": s.Description,
			"activation":  s.Activation,
			"expiration":  s.Expiration,
		}
	}

	return secrets
}

func resourceDataToOpenIdClient(d *schema.ResourceData) (*smilecdr.OpenIdClient, error) {

	secrets := d.Get("client_secrets").([]interface{})

	clientSecrets := []smilecdr.ClientSecret{}
	for _, secret := range secrets {
		s := secret.(map[string]interface{})
		if s["secret"] != nil || s["secret"].(string) != "" {
			secret := smilecdr.ClientSecret{
				Secret:      s["secret"].(string),
				Description: s["description"].(string),
				Activation:  s["activation"].(string),
				Expiration:  s["expiration"].(string),
			}
			clientSecrets = append(clientSecrets, secret)
		}
	}

	perms := d.Get("permissions").(*schema.Set).List()
	userPermissions := []smilecdr.UserPermission{}
	for _, perm := range perms {
		p := perm.(map[string]interface{})
		userPermissions = append(userPermissions, smilecdr.UserPermission{
			Permission: p["permission"].(string),
			Argument:   p["argument"].(string),
		})
	}

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
		NodeId:                      d.Get("node_id").(string),
		ModuleId:                    d.Get("module_id").(string),
		ClientId:                    d.Get("client_id").(string),
		ClientName:                  d.Get("client_name").(string),
		Enabled:                     d.Get("enabled").(bool),
		AccessTokenValiditySeconds:  d.Get("access_token_validity_seconds").(int),
		AllowedGrantTypes:           allowedGrantTypes,
		AlwaysRequireApproval:       d.Get("always_require_approval").(bool),
		AttestationAccepted:         d.Get("attestation_accepted").(bool),
		AutoApproveScopes:           autoApproveScopes,
		AutoGrantScopes:             autoGrantScopes,
		CanIntrospectAnyTokens:      d.Get("can_introspect_any_tokens").(bool),
		CanIntrospectOwnTokens:      d.Get("can_introspect_own_tokens").(bool),
		CanReissueTokens:            d.Get("can_reissue_tokens").(bool),
		ClientSecrets:               clientSecrets,
		CreatedByAppSphere:          d.Get("created_by_app_sphere").(bool),
		FixedScope:                  d.Get("fixed_scope").(bool),
		JwksUrl:                     d.Get("jwks_url").(string),
		Permissions:                 userPermissions,
		PublicJwks:                  d.Get("public_jwks").(string),
		RefreshTokenValiditySeconds: d.Get("refresh_token_validity_seconds").(int),
		RegisteredRedirectUris:      registeredRedirectUris,
		RememberApprovedScopes:      d.Get("remember_approved_scopes").(bool),
		Scopes:                      scopes,
		SecretClientCanChange:       d.Get("secret_client_can_change").(bool),
		SecretRequired:              d.Get("secret_required").(bool),
		ArchivedAt:                  d.Get("archived_at").(string),
	}

	return openidClient, nil

}

func resourceOpenIdClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	d.Set("created_by_app_sphere", false) // If we create a client, it is not created by AppSphere

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
	d.Set("node_id", openIdClient.NodeId)
	d.Set("module_id", openIdClient.ModuleId)
	d.Set("client_id", openIdClient.ClientId)
	d.Set("client_name", openIdClient.ClientName)
	d.Set("enabled", openIdClient.Enabled)
	d.Set("access_token_validity_seconds", openIdClient.AccessTokenValiditySeconds)
	d.Set("allowed_grant_types", openIdClient.AllowedGrantTypes)
	d.Set("always_require_approval", openIdClient.AlwaysRequireApproval)
	d.Set("attestation_accepted", openIdClient.AttestationAccepted)
	d.Set("auto_approve_scopes", openIdClient.AutoApproveScopes)
	d.Set("auto_grant_scopes", openIdClient.AutoGrantScopes)
	d.Set("can_introspect_any_tokens", openIdClient.CanIntrospectAnyTokens)
	d.Set("can_introspect_own_tokens", openIdClient.CanIntrospectOwnTokens)
	d.Set("can_reissue_tokens", openIdClient.CanReissueTokens)
	d.Set("client_secrets", flattenClientSecrets(openIdClient.ClientSecrets))
	d.Set("created_by_app_sphere", openIdClient.CreatedByAppSphere)
	d.Set("fixed_scope", openIdClient.FixedScope)
	d.Set("jwks_url", openIdClient.JwksUrl)
	d.Set("permissions", openIdClient.Permissions)
	d.Set("public_jwks", openIdClient.PublicJwks)
	d.Set("refresh_token_validity_seconds", openIdClient.RefreshTokenValiditySeconds)
	d.Set("registered_redirect_uris", openIdClient.RegisteredRedirectUris)
	d.Set("remember_approved_scopes", openIdClient.RememberApprovedScopes)
	d.Set("scopes", openIdClient.Scopes)
	d.Set("secret_client_can_change", openIdClient.SecretClientCanChange)
	d.Set("secret_required", openIdClient.SecretRequired)
	d.Set("archived_at", openIdClient.ArchivedAt)
	return diags

}

func resourceOpenIdClientUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	client, mErr := resourceDataToOpenIdClient(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}

	d.SetId(client.ClientId)

	if d.HasChange("client_secrets") {
		_, newValue := d.GetChange("client_secrets")
		if newValue.(*schema.Set).Len() == 0 {
			fmt.Println("Removing client secret")
			client.ClientSecrets = []smilecdr.ClientSecret{}
		}
	}

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
