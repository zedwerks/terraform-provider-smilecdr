// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
		Importer: &schema.ResourceImporter{
			StateContext: resourceOpenIdClientImport,
		},
		Schema: map[string]*schema.Schema{
			"created": {
				Type:     schema.TypeBool,
				Computed: true,
			},
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"secret": {
							Type:             schema.TypeString,
							Required:         true,
							Sensitive:        false,
							ValidateDiagFunc: validations.ValidateDiagFunc(validation.StringLenBetween(8, 256)),
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								created := d.Get("created").(bool) || false
								// Suppress the output of changes to the 'secret' attribute in the plan
								if created {
									fmt.Printf("k = %s, old = %s, new = %s\n", k, old, new)
									fmt.Println("Suppressing diff for sensitive attribute 'secret'")
								}
								return created
							},
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
			},
			"scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"secret_client_can_change": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"secret_required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"archived_at": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsRFC3339Time),
			},
		},
	}
}

func flattenClientSecrets(clientSecrets []smilecdr.ClientSecret) []interface{} {
	secrets := make([]interface{}, len(clientSecrets))

	for i, s := range clientSecrets {
		secrets[i] = map[string]interface{}{
			"pid":         s.Pid,
			"secret":      s.Secret,
			"description": s.Description,
			"activation":  s.Activation,
			"expiration":  s.Expiration,
		}
	}

	return secrets
}

func flattenPermissions(permissions []smilecdr.UserPermission) []interface{} {
	perms := make([]interface{}, len(permissions))

	for i, p := range permissions {
		perms[i] = map[string]interface{}{
			"permission": p.Permission,
			"argument":   p.Argument,
		}
	}

	return perms
}

func resourceDataToOpenIdClient(d *schema.ResourceData) (*smilecdr.OpenIdClient, error) {

	fmt.Println("In resourceDataToOpenIdClient")

	secrets := d.Get("client_secrets").(*schema.Set).List()

	clientSecrets := []smilecdr.ClientSecret{}
	for _, secret := range secrets {
		s := secret.(map[string]interface{})
		if s["secret"] != nil {
			secret := smilecdr.ClientSecret{
				Pid:         s["pid"].(int),
				Secret:      s["secret"].(string),
				Description: s["description"].(string),
				Activation:  s["activation"].(string),
				Expiration:  s["expiration"].(string),
			}
			clientSecrets = append(clientSecrets, secret)
		}
	}

	permissions := d.Get("permissions").(*schema.Set).List()

	userPermissions := make([]smilecdr.UserPermission, len(permissions))
	for _, permission := range permissions {
		s := permission.(map[string]interface{})
		if s["permission"] != nil || s["permission"].(string) != "" {
			perm := smilecdr.UserPermission{
				Permission: s["permission"].(string),
				Argument:   s["argument"].(string),
			}
			userPermissions = append(userPermissions, perm)
		}
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

	fmt.Println("In resourceOpenIdClientCreate")

	c := m.(*smilecdr.Client)

	d.Set("created_by_app_sphere", false) // If we create a client, it is not created by AppSphere
	d.Set("created", true)

	client, mErr := resourceDataToOpenIdClient(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}

	o, err := c.PostOpenIdClient(ctx, *client)

	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error creating openid client",
		})
		return diags
	}

	d.Set("created", true)   // Set the 'created' state variable to true after the initial creation
	d.SetId(client.ClientId) // the primary resource identifier. must be unique.
	d.Set("pid", o.Pid)      // the pid is needed for Put requests

	return resourceOpenIdClientRead(ctx, d, m)
}

func resourceOpenIdClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	fmt.Println("In resourceOpenIdClientRead")

	var diags diag.Diagnostics

	c := m.(*smilecdr.Client)

	client_id := d.Get("client_id").(string)
	nodeId := d.Get("node_id").(string)
	moduleId := d.Get("module_id").(string)

	openIdClient, err := c.GetOpenIdClient(ctx, nodeId, moduleId, client_id)

	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error reading openid client",
		})
		return diags
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
	d.Set("permissions", flattenPermissions(openIdClient.Permissions))
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

	fmt.Println("Updating OpenIdClient")

	c := m.(*smilecdr.Client)

	client, mErr := resourceDataToOpenIdClient(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}

	d.SetId(client.ClientId)

	_, err := c.PutOpenIdClient(ctx, *client)

	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error updating openid client",
		})
		return diags
	}

	return resourceOpenIdClientRead(ctx, d, m)

}

func resourceOpenIdClientDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	fmt.Println("Deleting OpenIdClient")

	var diags diag.Diagnostics

	d.Set("archived_at", time.Now().Format(time.RFC3339))

	c := m.(*smilecdr.Client)

	client, mErr := resourceDataToOpenIdClient(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}

	d.SetId(client.ClientId)

	_, err := c.PutOpenIdClient(ctx, *client)

	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error updating openid client",
		})
		return diags
	}

	d.SetId("")

	return diags
}

func resourceOpenIdClientImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(*smilecdr.Client)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import. supported import formats: {{nodeId}}/{{moduleId}}/{{clientId}}")
	}

	_, err := c.GetOpenIdClient(ctx, parts[0], parts[1], parts[3])
	if err != nil {
		return nil, err
	}

	d.Set("node_id", parts[0])
	d.Set("module_id", parts[1])
	d.Set("client_id", parts[2])

	diagnostics := resourceOpenIdClientRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
