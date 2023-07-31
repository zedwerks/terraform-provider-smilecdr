// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func resourceSmartOutboundSecurityCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceModuleConfigCreate(ctx, d, m)
}

func resourceSmartOutboundSecurityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceModuleConfigRead(ctx, d, m)
}

func resourceSmartOutboundSecurityUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceModuleConfigUpdate(ctx, d, m)
}

func resourceSmartOutboundSecurityDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceModuleConfigDelete(ctx, d, m)
}
