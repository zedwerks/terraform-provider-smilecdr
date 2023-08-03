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

func resourceSmartOutboundSecurity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSmartOutboundSecurityCreate,
		ReadContext:   resourceSmartOutboundSecurityRead,
		UpdateContext: resourceSmartOutboundSecurityUpdate,
		DeleteContext: resourceSmartOutboundSecurityDelete,
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
			"smart_capabilities_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of Smart Capabilities to enable (See http://hl7.org/fhir/smart-app-launch/conformance.html#capability-sets); one capability per line.",
			},
			"pkce_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If this setting is enabled, the server will require the use of PKCE for all Authorization Code SMART Auth flows. Enabling this setting also disallows the use of the OAuth2 Implicit Grant type, since this flow does not support PKCE.",
			},
			"pkce_plain_challenge_supported": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If this setting is enabled, the server will allow the use of the plain PKCE challenge method. This is not recommended, but is supported for backwards compatibility.",
			},
			"enforce_approved_scopes_to_restrict_permissions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"smart_configuration_scopes_supported": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "openid fhirUser launch launch/patient patient/*.* offline_access",
				Description: "",
			},
			"allowed_audience_list": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"federate_mode_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"davinci_consent_handling": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  9200,
			},
			"issuer_url": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.IsURLWithHTTPorHTTPS),
			},
			"openid_signing_keystore_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cors_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

func outboundSecurityResourceToModuleConfig(d *schema.ResourceData) (*smilecdr.ModuleConfig, error) {

	d.Set("module_type", "SECURITY_OUT_SMART") // Hardcoded for this module type

	moduleConfig := &smilecdr.ModuleConfig{
		ModuleId:   d.Get("module_id").(string),
		ModuleType: d.Get("module_type").(string),
		NodeId:     d.Get("node_id").(string),
	}

	if v, ok := d.GetOk("smart_capabilities_list"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "smart_capabilities_list",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("pkce_required"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "pkce.required",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("pkce_plain_challenge_supported"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "pkce.plain_challenge_supported",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("enforce_approved_scopes_to_restrict_permissions"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "enforce_approved_scopes_to_restrict_permissions",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("smart_configuration_scopes_supported"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "smart_configuration.scopes_supported",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("allowed_audience_list"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "allowed_audience_list",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("federate_mode_enabled"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "federate_mode.enabled",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("davinci_consent_handling"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "davinci.consent_handling",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("port"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "port",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("issuer_url"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "issuer.url",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("openid_signing_keystore_id"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "openid.signing.keystore_id",
			Value: v.(string),
		})
	}

	if v, ok := d.GetOk("cors_enable"); ok {
		moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
			Key:   "cors.enable",
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

func resourceSmartOutboundSecurityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	moduleConfig, err := outboundSecurityResourceToModuleConfig(d)

	if err != nil {
		return diag.FromErr(err)
	}

	module, err := c.PostModuleConfig(*moduleConfig)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(module.ModuleId) // the primary resource identifier. must be unique.

	return resourceSmartOutboundSecurityRead(ctx, d, m)

}

func resourceSmartOutboundSecurityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	moduleId := d.Get("module_id").(string)
	nodeId := d.Get("node_id").(string)
	moduleConfig, err := c.GetModuleConfig(nodeId, moduleId)

	// map from moduleConfig to resourceData
	if err != nil {
		return diag.FromErr(err)
	}

	val, ok := moduleConfig.LookupOptionOk("smart_capabilities_list")
	if ok {
		d.Set("smart_capabilities_list", val)
	}

	val, ok = moduleConfig.LookupOptionOk("pkce.required")
	if ok {
		d.Set("pkce_required", val)
	}

	val, ok = moduleConfig.LookupOptionOk("pkce.plain_challenge_supported")
	if ok {
		d.Set("pkce_plain_challenge_supported", val)
	}

	val, ok = moduleConfig.LookupOptionOk("enforce_approved_scopes_to_restrict_permissions")
	if ok {
		d.Set("enforce_approved_scopes_to_restrict_permissions", val)
	}

	val, ok = moduleConfig.LookupOptionOk("smart_configuration.scopes_supported")
	if ok {
		d.Set("smart_configuration_scopes_supported", val)
	}

	val, ok = moduleConfig.LookupOptionOk("allowed_audience_list")
	if ok {
		d.Set("allowed_audience_list", val)
	}

	val, ok = moduleConfig.LookupOptionOk("federate_mode.enabled")
	if ok {
		d.Set("federate_mode_enabled", val)
	}

	val, ok = moduleConfig.LookupOptionOk("davinci.consent_handling")
	if ok {
		d.Set("davinci_consent_handling", val)
	}

	val, ok = moduleConfig.LookupOptionOk("port")
	if ok {
		d.Set("port", val)
	}

	val, ok = moduleConfig.LookupOptionOk("issuer.url")
	if ok {
		d.Set("issuer_url", val)
	}

	val, ok = moduleConfig.LookupOptionOk("openid.signing.keystore_id")
	if ok {
		d.Set("openid_signing_keystore_id", val)
	}

	val, ok = moduleConfig.LookupOptionOk("cors.enable")
	if ok {
		d.Set("cors_enable", val)
	}

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

func resourceSmartOutboundSecurityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	moduleConfig, err := outboundSecurityResourceToModuleConfig(d)

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(moduleConfig.ModuleId) // the primary resource identifier. must be unique.

	_, pErr := c.PutModuleConfig(*moduleConfig)

	if pErr != nil {
		return diag.FromErr(pErr)
	}

	return resourceSmartOutboundSecurityRead(ctx, d, m)
}

func resourceSmartOutboundSecurityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
