// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"smart_capabilities_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				Description: "List of Smart Capabilities to enable (See http://hl7.org/fhir/smart-app-launch/conformance.html#capability-sets); one capability per line.",
			},
			"pkce.required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If this setting is enabled, the server will require the use of PKCE for all Authorization Code SMART Auth flows. Enabling this setting also disallows the use of the OAuth2 Implicit Grant type, since this flow does not support PKCE.",
			},
			"pkce.plain_challenge_supported": {
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
			"smart_configuration.scopes_supported": {
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
			"federate_mode.enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"davinci.consent_handling.enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  9200,
			},
			"issuer.url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"openid.signing.keystore_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cors.enable": {
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

	moduleConfig := &smilecdr.ModuleConfig{
		ModuleId:   d.Get("module_id").(string),
		ModuleType: d.Get("module_type").(string),
	}

	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "smart_capabilities_list",
		Value: d.Get("smart_capabilities_list").(string),
	})
	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "pkce.required",
		Value: d.Get("pkce.required").(string),
	})
	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "pkce.plain_challenge_supported",
		Value: d.Get("pkce.plain_challenge_supported").(string),
	})
	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "enforce_approved_scopes_to_restrict_permissions",
		Value: d.Get("enforce_approved_scopes_to_restrict_permissions").(string),
	})
	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "smart_configuration.scopes_supported",
		Value: d.Get("smart_configuration.scopes_supported").(string),
	})
	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "allowed_audience_list",
		Value: d.Get("allowed_audience_list").(string),
	})
	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "federate_mode.enabled",
		Value: d.Get("federate_mode.enabled").(string),
	})
	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "davinci.consent_handling.enabled",
		Value: d.Get("davinci.consent_handling.enabled").(string),
	})
	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "port",
		Value: d.Get("port").(string),
	})
	moduleConfig.Options = append(moduleConfig.Options, smilecdr.ModuleOption{
		Key:   "issuer.url",
		Value: d.Get("issuer.url").(string),
	})

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
	moduleConfig, err := c.GetModuleConfig(moduleId)

	// map from moduleConfig to resourceData
	if err != nil {
		return diag.FromErr(err)
	}
	val, err := moduleConfig.LookupOption(("smart_capabilities_list"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("smart_capabilities_list", val)
	val, err = moduleConfig.LookupOption(("pkce.required"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("pkce.required", val)
	val, err = moduleConfig.LookupOption(("pkce.plain_challenge_supported"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("pkce.plain_challenge_supported", val)
	val, err = moduleConfig.LookupOption(("enforce_approved_scopes_to_restrict_permissions"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("enforce_approved_scopes_to_restrict_permissions", val)
	val, err = moduleConfig.LookupOption(("smart_configuration.scopes_supported"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("smart_configuration.scopes_supported", val)
	val, err = moduleConfig.LookupOption(("allowed_audience_list"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("allowed_audience_list", val)
	val, err = moduleConfig.LookupOption(("federate_mode.enabled"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("federate_mode.enabled", val)
	val, err = moduleConfig.LookupOption(("davinci.consent_handling.enabled"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("davinci.consent_handling.enabled", val)
	val, err = moduleConfig.LookupOption(("port"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("port", val)
	val, err = moduleConfig.LookupOption(("issuer.url"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("issuer.url", val)
	val, err = moduleConfig.LookupOption(("openid.signing.keystore_id"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("openid.signing.keystore_id", val)
	val, err = moduleConfig.LookupOption(("cors.enable"))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("cors.enable", val)

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

	err := c.DeleteModuleConfig(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("") // This is unset when the resource is deleted

	return nil
}
