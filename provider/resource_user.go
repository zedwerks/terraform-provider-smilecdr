// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zedwerks/terraform-smilecdr/provider/helper/validations"
	"github.com/zedwerks/terraform-smilecdr/smilecdr"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImport,
		},
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
				Default:  "local_security",
			},
			"username": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"password": {
				Type:             schema.TypeString,
				Required:         true,
				Sensitive:        true,
				DiffSuppressFunc: smilecdr.suppressSensitiveDataDiff,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.StringLenBetween(8, 512)),
			},
			"family_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"given_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_active": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_connected": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_locked": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_user": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"account_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"external": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"service_account": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"2fa_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorities": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission": {
							Type:     schema.TypeString,
							Required: true,
						},
						"argument": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func flattenUserAuthorities(permissions []smilecdr.UserAuthorities) []interface{} {
	perms := make([]interface{}, len(permissions))

	for i, p := range permissions {
		perms[i] = map[string]interface{}{
			"permission": p.Permission,
			"argument":   p.Argument,
		}
	}

	return perms
}

func resourceDataToUser(d *schema.ResourceData) (*smilecdr.User, error) {

	fmt.Println("In resourceDataToUser...")

	//authorities := d.Get("authorities").([]interface{})

	authorities := d.Get("authorities").(*schema.Set).List()

	userAuthorities := []smilecdr.UserAuthorities{}
	for _, authority := range authorities {
		s := authority.(map[string]interface{})
		if s["permission"] != nil || s["permission"].(string) != "" {
			perm := smilecdr.UserAuthorities{
				Permission: s["permission"].(string),
				Argument:   s["argument"].(string),
			}
			userAuthorities = append(userAuthorities, perm)
		}
	}

	smileUser := &smilecdr.User{
		Pid:                 d.Get("pid").(int),
		NodeId:              d.Get("node_id").(string),
		ModuleId:            d.Get("module_id").(string),
		Username:            d.Get("username").(string),
		Password:            d.Get("password").(string),
		FamilyName:          d.Get("family_name").(string),
		GivenName:           d.Get("given_name").(string),
		AccountLocked:       d.Get("account_locked").(bool),
		SystemUser:          d.Get("system_user").(bool),
		AccountDisabled:     d.Get("account_disabled").(bool),
		External:            d.Get("external").(bool),
		ServiceAccount:      d.Get("service_account").(bool),
		TwoFactorAuthStatus: d.Get("2fa_status").(string),
		Authorities:         userAuthorities,
	}

	return smileUser, nil

}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	user, mErr := resourceDataToUser(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}

	o, err := c.PostUser(ctx, *user)
	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error creating openid client",
		})
		return diags
	}
	d.Set("pid", o.Pid)
	d.SetId(strconv.Itoa(o.Pid)) // the primary resource identifier. must be unique.

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	fmt.Println("Read User")

	var diags diag.Diagnostics

	c := m.(*smilecdr.Client)

	pid := d.Get("pid").(int)
	nodeId := d.Get("node_id").(string)
	moduleId := d.Get("module_id").(string)

	user, err := c.GetUser(ctx, nodeId, moduleId, pid)

	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error reading user record",
		})
		return diags
	}

	d.SetId(strconv.Itoa(user.Pid))
	d.Set("pid", user.Pid)
	d.Set("node_id", user.NodeId)
	d.Set("module_id", user.ModuleId)
	d.Set("username", user.Username)
	d.Set("password", "") // Sensitive Data not returned by Smile CDR
	d.Set("family_name", user.FamilyName)
	d.Set("given_name", user.GivenName)
	d.Set("account_locked", user.AccountLocked)
	d.Set("system_user", user.SystemUser)
	d.Set("authorities", flattenUserAuthorities(user.Authorities))
	d.Set("account_disabled", user.AccountDisabled)
	d.Set("external", user.External)
	d.Set("service_account", user.ServiceAccount)
	d.Set("2fa_status", user.TwoFactorAuthStatus)

	return diags

}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	fmt.Println("Updating User")

	c := m.(*smilecdr.Client)

	user, mErr := resourceDataToUser(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}

	d.SetId(strconv.Itoa(user.Pid))

	_, err := c.PutUser(ctx, *user)

	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error updating user record",
		})
		return diags
	}

	return resourceUserRead(ctx, d, m)

}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	fmt.Println("Deleting User")

	c := m.(*smilecdr.Client)

	var diags diag.Diagnostics

	d.Set("account_disabled", true) // disable the account

	user, mErr := resourceDataToUser(d)
	if mErr != nil {
		return diag.FromErr(mErr)
	}
	d.SetId(strconv.Itoa(user.Pid))

	_, err := c.PutUser(ctx, *user)

	if err != nil {
		diags := diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Error updating user record to disable account",
		})
		return diags
	}

	return diags
}

func resourceUserImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	c := meta.(*smilecdr.Client)

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import. supported import formats: {{nodeId}}/{{moduleId}}/{{pid}}")
	}

	pid, _ := strconv.Atoi(parts[2])

	_, err := c.GetUser(ctx, parts[0], parts[1], pid)
	if err != nil {
		return nil, err
	}

	d.Set("node_id", parts[0])
	d.Set("module_id", parts[1])
	d.Set("pid", parts[2])
	d.SetId(parts[2])

	diagnostics := resourceUserRead(ctx, d, meta)
	if diagnostics.HasError() {
		return nil, errors.New(diagnostics[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
