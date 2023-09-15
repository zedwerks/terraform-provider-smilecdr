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

func resourceModuleConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModuleConfigCreate,
		ReadContext:   resourceModuleConfigRead,
		UpdateContext: resourceModuleConfigUpdate,
		DeleteContext: resourceModuleConfigDelete,
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Master",
			},
			"module_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"module_type": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateDiagFunc(validation.StringInSlice([]string{"LICENSE", "SECURITY_IN_LOCAL", "SECURITY_OUT_SMART"}, false)),
			},
			"options": {
				Type:     schema.TypeList,
				Required: true,
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

func resourceToModuleConfig(d *schema.ResourceData) (*smilecdr.ModuleConfig, error) {
	moduleConfig := &smilecdr.ModuleConfig{
		ModuleId:   d.Get("module_id").(string),
		ModuleType: d.Get("module_type").(string),
	}

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
		moduleConfig.Dependencies = append(moduleConfig.Dependencies, smilecdr.ModuleDependency{
			ModuleId: dependencyMap["module_id"].(string),
			Type:     dependencyMap["type"].(string),
		})
	}

	return moduleConfig, nil
}

func resourceModuleConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*smilecdr.Client)

	moduleConfig, mErr := resourceToModuleConfig(d)
	nodeId := d.Get("node_id").(string)

	if mErr != nil {
		return diag.FromErr(mErr)
	}

	_, err := c.PostModuleConfig(nodeId, *moduleConfig)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(moduleConfig.ModuleId) // the primary resource identifier. must be unique.

	return resourceModuleConfigRead(ctx, d, m)
}

func resourceModuleConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*smilecdr.Client)

	moduleId := d.Get("module_id").(string)
	nodeId := d.Get("node_id").(string)

	moduleConfig, err := c.GetModuleConfig(nodeId, moduleId)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("module_id", moduleConfig.ModuleId)
	d.Set("module_type", moduleConfig.ModuleType)

	options := make([]interface{}, len(moduleConfig.Options))
	for i, option := range moduleConfig.Options {
		options[i] = map[string]interface{}{
			"key":   option.Key,
			"value": option.Value,
		}
	}
	d.Set("options", options)

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

func resourceModuleConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*smilecdr.Client)

	moduleConfig, mErr := resourceToModuleConfig(d)
	nodeId := d.Get("node_id").(string)

	if mErr != nil {
		return diag.FromErr(mErr)
	}

	d.SetId(moduleConfig.ModuleId)

	_, err := c.PutModuleConfig(nodeId, *moduleConfig)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceModuleConfigRead(ctx, d, m)
}

func resourceModuleConfigDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
