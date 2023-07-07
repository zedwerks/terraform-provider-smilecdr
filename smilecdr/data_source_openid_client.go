// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpenIdClient() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceOpenIdClientRead,
		Schema: map[string]*schema.Schema{
			"pid": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
			},
			"nodeId": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "Master",
			},
			"moduleId": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "smart_auth",
			},
			"accessTokenValiditySeconds": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				Default:  300,
			},
			"allowedGrantTypes": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Schema{
					Type:    schema.TypeString,
				},
			},
			"autoApproveScopes": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
				},
			},
			"autoGrantScopes": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
				},
			},
			"clientId": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"clientName": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "Some Client",
			},
			"clientSecrets": &schema.Schema{
				Type:     schema.TypeList,
				Required: false,
				Elem: &schema.Schema{
						"pid": &schema.Schema{
							Type:     schema.TypeInt,
							Required: false,
						},
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
							Type:     schema.TypeString,
							Required: false,
							Default:  "",
							ValidateFunc: validation.IsRFC3339Time,
						},
						"expiration": &schema.Schema{
							Type:     schema.TypeString,
							Required: false,
							Default:  "",
							ValidateFunc: validation.IsRFC3339Time,
						},
					},
				},
				"fixedScope": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
				"refreshTokenValiditySeconds": &schema.Schema{
					Type:     schema.TypeInt,
					Required: false,
					Default:  86400,
				},
				"registeredRedirectUris": &schema.Schema{
					Type:     schema.TypeList,
					Required: false,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"scopes": &schema.Schema{
					Type:     schema.TypeList,
					Required: false,
					Elem: &schema.Schema{
						Type:     schema.TypeString,
					},
				},
				"secretRequired": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
				"secretClientCanChange": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
				"enabled": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  true,
				},
				"canIntrospectAnyTokens": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
				"canIntrospectOwnTokens": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
				"alwaysRequireApproval": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
				"canReissueTokens": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
				"permissions": &schema.Schema{
					Type:     schema.TypeList,
					Required: false,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"rememberApprovedScopes": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
				"attestationAccepted": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
				"publicJwks": &schema.Schema{
					Type:     schema.TypeString,
					Required: false,
					Default:  "",
				},
				"jwksUrl": &schema.Schema{
					Type:     schema.TypeString,
					Required: false,
					Default:  "",
				},
				"archivedAt": &schema.Schema{
					Type:     schema.TypeString,
					Required: false,
					ValidateFunc: validation.IsRFC3339Time,
				},
				"createdByAppSphere": &schema.Schema{
					Type:     schema.TypeBool,
					Required: false,
					Default:  false,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceOpenIdClientRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	return diags

}

