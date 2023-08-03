// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zedwerks/terraform-smilecdr/provider/helper/validations"
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
			},
			"node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Master",
				Description: "The node ID of the node to be configured.",
			},

			"enforce_approved_scopes_to_restrict_permissions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, only scopes that have been approved for the client will be used to determine the permissions that the client has. If false, all scopes that are associated with the client will be used to determine the permissions that the client has.",
			},
			"trust_intra_cluster_tokens_modules": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A list of module IDs that are trusted to issue tokens that are valid for intra-cluster communication. If a token is received from a module that is not in this list, it will be rejected.",
			},
			"cache_authentication_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Successfully validated authentication tokens will be cached for the given number of seconds. This cache has a positive impact on performance since validating a token is a non-trivial operation and may involve lookups and crypto operations. However, setting this timeout to a value that is high means that invalidated tokens will be trusted for longer than they should be",
			},
			"key_validation_prevent_token_key_reuse": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the same key will not be used to sign multiple tokens. This is a security measure that prevents a key that has been compromised from being used to sign new tokens.",
			},
			"key_validation_require_key_expiry": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, tokens will only be accepted if they are signed with a key that has an expiry date. This is a security measure that prevents a key that has been compromised from being used to sign new tokens.",
			},
			"smart_configuration_scopes_supported": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "openid profile email",
				Description: "A space-separated list of scopes that are supported by the SMART on FHIR server. This list is used to validate the scopes that are requested by the client. If the client requests a scope that is not in this list, the request will be rejected.",
			},
			"token_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
				Description:      "The URL of the token endpoint. This is the endpoint that the SMART on FHIR client will use to obtain an access token.",
			},
			"authorization_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
				Description:      "The URL of the authorization endpoint. This is the endpoint that the SMART on FHIR client will use to obtain an authorization code.",
			},
			"management_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
				Description:      "The URL of the management endpoint. This is the endpoint that the SMART on FHIR client will use to obtain a refresh token.",
			},
			"introspection_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
				Description:      "The URL of the introspection endpoint. This is the endpoint that the SMART on FHIR client will use to validate an access token.",
			},
			"revocation_endpoint": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
				Description:      "The URL of the revocation endpoint. This is the endpoint that the SMART on FHIR client will use to revoke an access token.",
			},
			"introspection_client_jwks_cache_mins": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "The minutes the keystore is valid.  If set to a non-zero value, any keystore lookups performed by the OIDC HTTP Client will be cached for the specified number of minutes. Caching these fetched keystores improves authentication performance by avoiding unnecessary lookups, but can also mean that invalidated keys will be honored for a period. Setting this to a small setting (such as the default value) is generally a sensible compromise.",
			},
			"introspection_client_truststore_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to the trust store file. If set, the trust store file will be used to validate the TLS certificate of the introspection endpoint. If not set, the introspection endpoint will not be validated.",
			},
			"callback_script_text": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The text of the callback script. This script will be executed when the SMART on FHIR client has been successfully authenticated. The script will be executed in the context of the authenticated user. The script can be used to perform custom actions when the user has been authenticated.",
			},
			"tfa_totp_issuer_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The issuer name that will be used when generating TOTP tokens. This name will be displayed to the user when they are configuring their TOTP client.",
			},
			"tfa_totp_lock_after_failed_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"seed_servers_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to the seed servers file. This file contains a list of seed servers that will be used to bootstrap the cluster. If this file is not set, the node will not be able to join the cluster.",
			},
			"debug_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"debug_host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"debug_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"debug_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"dependencies": {
				Type:     schema.TypeList,
				Optional: true,
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

	if v, ok := d.GetOk("enforce_approved_scopes_to_restrict_permissions"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "enforce_approved_scopes_to_restrict_permissions",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("trust_intra_cluster_tokens_modules"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "trust_intra_cluster_tokens.modules",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("cache_authentication_seconds"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "cache_authentication.seconds",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("key_validation_prevent_token_key_reuse"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "key_validation.prevent_token_key_reuse",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("key_validation_require_key_expiry"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "key_validation.require_key_expiry",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("smart_configuration_scopes_supported"); ok {
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
	if v, ok := d.GetOk("management_endpoint"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "management_endpoint",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("introspection_endpoint"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "introspection_endpoint",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("revocation_endpoint"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "revocation_endpoint",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("introspection_client_jwks_cache_mins"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "introspection_client.jwks_cache.mins",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("introspection_client_truststore_file"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "introspection_client.truststore.file",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("callback_script_text"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "callback_script.text",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tfa_totp_issuer_name"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tfa.totp.issuer_name",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("tfa_totp_lock_after_failed_attempts"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "tfa.totp.lock_after_failed_attempts",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("seed_servers_file"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "seed_servers.file",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("debug_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.debug_enabled",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("debug_host"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.host",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("debug_port"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.port",
			Value: v.(string),
		})
	}
	if v, ok := d.GetOk("debug_path"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "debug.path",
			Value: v.(string),
		})
	}

	// Add any other options that are not in the schema
	options := d.Get("options").([]interface{})
	for _, option := range options {
		optionMap := option.(map[string]interface{})
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   optionMap["key"].(string),
			Value: optionMap["value"].(string),
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
	c := m.(*smilecdr.Client)

	moduleId := d.Get("module_id").(string)
	nodeId := d.Get("node_id").(string)
	moduleConfig, err := c.GetModuleConfig(nodeId, moduleId)

	// map from moduleConfig to resourceData
	if err != nil {
		return diag.FromErr(err)
	}
	val, err := moduleConfig.LookupOption("enforce_approved_scopes_to_restrict_permissions")
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("enforce_approved_scopes_to_restrict_permissions", val)

	val, err = moduleConfig.LookupOption("trust_intra_cluster_tokens.modules")
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("trust_intra_cluster_tokens_modules", val)

	val, err = moduleConfig.LookupOption("cache_authentication.seconds")
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("cache_authentication_seconds", val)

	val, err = moduleConfig.LookupOption(("key_validation.prevent_token_key_reuse"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("key_validation_prevent_token_key_reuse", val)

	val, err = moduleConfig.LookupOption(("key_validation.require_key_expiry"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("key_validation_require_key_expiry", val)

	val, err = moduleConfig.LookupOption(("smart_configuration.scopes_supported"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("smart_configuration_scopes_supported", val)

	val, err = moduleConfig.LookupOption(("token_endpoint"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("token_endpoint", val)

	val, err = moduleConfig.LookupOption(("authorization_endpoint"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("authorization_endpoint", val)

	val, err = moduleConfig.LookupOption(("management_endpoint"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("management_endpoint", val)

	val, err = moduleConfig.LookupOption(("introspection_endpoint"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("introspection_endpoint", val)

	val, err = moduleConfig.LookupOption(("revocation_endpoint"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("revocation_endpoint", val)

	val, err = moduleConfig.LookupOption(("introspection_client.jwks_cache.mins"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("introspection_client_jwks_cache_mins", val)

	val, err = moduleConfig.LookupOption(("introspection_client.truststore.file"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("introspection_client_truststore_file", val)

	val, err = moduleConfig.LookupOption(("callback_script.text"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("callback_script_text", val)

	val, err = moduleConfig.LookupOption(("tfa.totp.issuer_name"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("tfa_totp_issuer_name", val)

	val, err = moduleConfig.LookupOption(("tfa.totp.lock_after_failed_attempts"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("tfa_totp_lock_after_failed_attempts", val)

	val, err = moduleConfig.LookupOption(("seed_servers.files"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("seed_servers_file", val)

	val, err = moduleConfig.LookupOption(("debug.debug_enabled"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("debug_enabled", val)

	val, err = moduleConfig.LookupOption(("debug.host"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("debug_host", val)

	val, err = moduleConfig.LookupOption(("debug.port"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("debug_port", val)

	val, err = moduleConfig.LookupOption(("debug.path"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("debug_path", val)

	options := make([]interface{}, len(moduleConfig.Options))
	for i, option := range moduleConfig.Options {
		options[i] = map[string]interface{}{
			"key":   option.Key,
			"value": option.Value,
		}
	}
	d.Set("options", options)

	// Set The Dependencies
	dependencies := make([]interface{}, len(moduleConfig.Dependencies))
	for i, dependency := range moduleConfig.Dependencies {
		dependencies[i] = map[string]interface{}{
			"module_id": dependency.ModuleId,
			"type":      dependency.Type,
		}
	}
	d.Set("dependencies", dependencies)

	return nil
}

func resourceSmartInboundSecurityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*smilecdr.Client)

	moduleConfig, err := inboundSecurityResourceToModuleConfig(d)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(moduleConfig.ModuleId) // the primary resource identifier. must be unique.

	_, pErr := c.PutModuleConfig(*moduleConfig)

	if pErr != nil {
		return diag.FromErr(pErr)
	}

	return resourceSmartInboundSecurityRead(ctx, d, m)
}

func resourceSmartInboundSecurityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*smilecdr.Client)

	moduleId := d.Get("module_id").(string)
	nodeId := d.Get("node_id").(string)

	err := c.DeleteModuleConfig(nodeId, moduleId)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("") // This is unset when the resource is deleted

	return nil
}
