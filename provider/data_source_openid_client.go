// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpenIdClients() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenIdClientsRead,
		Schema: map[string]*schema.Schema{
			"clients": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
							Default:  "Master",
						},
						"module_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Optional: true,
							Default:  "smart_auth",
						},
						"access_token_validity_seconds": &schema.Schema{
							Type:     schema.TypeInt,
							Required: false,
							Optional: true,
							Default:  300,
						},
						"allowed_grant_types": &schema.Schema{
							Type:     schema.TypeSet,
							Required: false,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								Required: false,
							},
						},
						"auto_approve_scopes": &schema.Schema{
							Type:     schema.TypeSet,
							Required: false,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"auto_grant_scopes": &schema.Schema{
							Type:     schema.TypeSet,
							Required: false,
							Optional: true,
							Default:  []string{""},
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"client_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"client_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							Optional: false,
						},
						"client_secrets": &schema.Schema{
							Type:     schema.TypeSet,
							Required: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"secret": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Required: false,
										Default:  "",
									},
									"activation": &schema.Schema{
										Type:         schema.TypeString,
										Required:     false,
										Default:      "",
										ValidateFunc: validation.IsRFC3339Time,
									},
									"expiration": &schema.Schema{
										Type:         schema.TypeString,
										Required:     false,
										Default:      "",
										ValidateFunc: validation.IsRFC3339Time,
									},
								},
							},
						},
						"fixed_scope": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
						"refresh_token_validity_seconds": &schema.Schema{
							Type:     schema.TypeInt,
							Required: false,
							Default:  86400,
						},
						"registered_redirect_uris": &schema.Schema{
							Type:     schema.TypeSet,
							Required: false,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scopes": &schema.Schema{
							Type:     schema.TypeSet,
							Required: false,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"secret_required": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
						"secret_client_can_change": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  true,
						},
						"can_introspect_any_tokens": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
						"can_introspect_own_tokens": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
						"always_require_approval": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
						"can_reissue_tokens": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
						"permissions": &schema.Schema{
							Type:     schema.TypeSet,
							Required: false,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"remember_approved_scopes": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
						"attestation_accepted": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
						"public_jwks": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Default:  "",
						},
						"jwks_url": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Default:  "",
						},
						"archivedAt": &schema.Schema{
							Type:         schema.TypeString,
							Required:     false,
							ValidateFunc: validation.IsRFC3339Time,
						},
						"created_by_app_sphere": &schema.Schema{
							Type:     schema.TypeBool,
							Required: false,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func dataSourceOpenIdClientsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
