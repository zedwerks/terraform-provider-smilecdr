// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zedwerks/terraform-smilecdr/smilecdr"
)

func resourceSmartInboundSecurity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSmartInboundSecurityCreate,
		ReadContext:   resourceSmartInboundSecurityRead,
		UpdateContext: resourceSmartInboundSecurityUpdate,
		DeleteContext: resourceSmartInboundSecurityDelete,
		Schema: map[string]*schema.Schema{
			"module_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique module ID of the module to be configured.",
			},
			"module_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The module type of the module to be configured.",
				Default:     "SECURITY_IN_SMART",
			},
			"node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Master",
				Description: "The node ID of the node to be configured.",
			},
			"cache_authentication.seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Successfully validated authentication tokens will be cached for the given number of seconds. This cache has a positive impact on performance since validating a token is a non-trivial operation and may involve lookups and crypto operations. However, setting this timeout to a value that is high means that invalidated tokens will be trusted for longer than they should be",
			},
			"key_validation.prevent_token_key_reuse": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the same key will not be used to sign multiple tokens. This is a security measure that prevents a key that has been compromised from being used to sign new tokens.",
			},
			"key_validation.require_key_expiry": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, tokens will only be accepted if they are signed with a key that has an expiry date. This is a security measure that prevents a key that has been compromised from being used to sign new tokens.",
			},
			"enforce_approved_scopes_to_restrict_permissions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, only scopes that have been approved for the client will be used to determine the permissions that the client has. If false, all scopes that are associated with the client will be used to determine the permissions that the client has.",
			},
			"smart_configuration.scopes_supported": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "openid profile email",
				Description: "A space-separated list of scopes that are supported by the SMART on FHIR server. This list is used to validate the scopes that are requested by the client. If the client requests a scope that is not in this list, the request will be rejected.",
			},
			"token_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the token endpoint. This is the endpoint that the SMART on FHIR client will use to obtain an access token.",
			},
			"authorization_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the authorization endpoint. This is the endpoint that the SMART on FHIR client will use to obtain an authorization code.",
			},
			"management_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the management endpoint. This is the endpoint that the SMART on FHIR client will use to obtain a refresh token.",
			},
			"introspection_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the introspection endpoint. This is the endpoint that the SMART on FHIR client will use to validate an access token.",
			},
			"revocation_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL of the revocation endpoint. This is the endpoint that the SMART on FHIR client will use to revoke an access token.",
			},
			"introspection_client.jwks_cache.mins": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "The minutes the keystore is valid.  If set to a non-zero value, any keystore lookups performed by the OIDC HTTP Client will be cached for the specified number of minutes. Caching these fetched keystores improves authentication performance by avoiding unnecessary lookups, but can also mean that invalidated keys will be honored for a period. Setting this to a small setting (such as the default value) is generally a sensible compromise.",
			},
			"introspection_client.truststore.file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to the trust store file. If set, the trust store file will be used to validate the TLS certificate of the introspection endpoint. If not set, the introspection endpoint will not be validated.",
			},
			"callback_script.text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The text of the callback script. This script will be executed when the SMART on FHIR client has been successfully authenticated. The script will be executed in the context of the authenticated user. The script can be used to perform custom actions when the user has been authenticated.",
			},
			"tfa.totp.issuer_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The issuer name that will be used when generating TOTP tokens. This name will be displayed to the user when they are configuring their TOTP client.",
			},
			"tfa.totp.lock_after_failed_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"seed_servers.file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to the seed servers file. This file contains a list of seed servers that will be used to bootstrap the cluster. If this file is not set, the node will not be able to join the cluster.",
			},
			"debug.debug_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"debug.host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"debug.port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"debug.path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dependencies": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"module_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func inboundSecurityResourceToModuleConfig(d *schema.ResourceData) (*smilecdr.ModuleConfig, error) {

	d.Set("module_type", "SECURITY_IN_SMART") // Hardcoded for this module type

	moduleConfig := &smilecdr.ModuleConfig{
		ModuleId:   d.Get("module_id").(string),
		ModuleType: d.Get("module_type").(string),
		NodeId:     d.Get("node_id").(string),
	}

	if v, ok := d.GetOk("cache_authentication.seconds"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "cache_authentication.seconds",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("key_validation.prevent_token_key_reuse"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "key_validation.prevent_token_key_reuse",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("key_validation.require_key_expiry"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "key_validation.require_key_expiry",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("enforce_approved_scopes_to_restrict_permissions"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "enforce_approved_scopes_to_restrict_permissions",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("smart_configuration.scopes_supported"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "smart_configuration.scopes_supported",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("token_endpoint"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "token_endpoint",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("authorization_endpoint"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "authorization_endpoint",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("introspection_endpoint"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "introspection_endpoint",
			Value: v.(string),
		})
	}

	dependencies := d.Get("dependencies").([]interface{})
	for _, dependency := range dependencies {
		dependencyMap := dependency.(map[string]interface{})
		moduleConfig.Dependencies = append(moduleConfig.Dependencies, smilecdr.ModuleDependencies{
			ModuleId: dependencyMap["module_id"].(string),
			Type:     dependencyMap["type"].(string),
		})
	}

	return moduleConfig, nil
}

func resourceSmartInboundSecurityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*smilecdr.Client)

	moduleConfig, err := inboundSecurityResourceToModuleConfig(d)

	if err != nil {
		return diag.FromErr(err)
	}

	module, err := c.PostModuleConfig(*moduleConfig)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(module.ModuleId) // the primary resource identifier. must be unique.

	return resourceSmartInboundSecurityRead(ctx, d, m)
}

func resourceSmartInboundSecurityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceSmartInboundSecurityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}

func resourceSmartInboundSecurityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	return diags
}
