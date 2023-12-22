// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"context"
	"encoding/json"
	"fmt"
)

type ModuleOption struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type ModuleDependency struct {
	ModuleId string `json:"moduleId,omitempty"`
	Type     string `json:"type,omitempty"`
}

type ModuleConfig struct {
	ModuleId     string             `json:"moduleId,omitempty"`
	ModuleType   string             `json:"moduleType,omitempty"`
	Options      []ModuleOption     `json:"options,omitempty"`
	Dependencies []ModuleDependency `json:"dependencies,omitempty"`
}

func printModuleConfig(moduleConfig ModuleConfig) {
	fmt.Println("ModuleConfig:")
	fmt.Println("  ModuleId: ", moduleConfig.ModuleId)
	fmt.Println("  ModuleType: ", moduleConfig.ModuleType)
	fmt.Println("  Options:")
	for _, kv := range moduleConfig.Options {
		fmt.Println("    ", kv.Key, ": ", kv.Value)
	}
	fmt.Println("  Dependencies:")
	for _, dep := range moduleConfig.Dependencies {
		fmt.Println("    ", dep.ModuleId, ": ", dep.Type)
	}
}

func (moduleConfig *ModuleConfig) LookupOptionOk(key string) (string, bool) {

	moduleConfigOptions := moduleConfig.Options
	for _, kv := range moduleConfigOptions {
		if kv.Key == key {
			return kv.Value, true
		}
	}

	return "", false
}

func (smilecdr *Client) GetModuleConfigs(ctx context.Context) ([]ModuleConfig, error) {
	var modules []ModuleConfig
	jsonBody, getErr := smilecdr.Get(ctx, "/module-config")
	if getErr != nil {
		return modules, getErr
	}

	err := json.Unmarshal(jsonBody, &modules)

	return modules, err
}

func (smilecdr *Client) GetModuleConfig(ctx context.Context, nodeId string, moduleId string) (ModuleConfig, error) {
	var module ModuleConfig
	var endpoint = fmt.Sprintf("/module-config/%s/%s", nodeId, moduleId)
	jsonBody, getErr := smilecdr.Get(ctx, endpoint)
	if getErr != nil {
		fmt.Println("error during Get in GetModuleConfig:", getErr)
		return module, getErr
	}

	err := json.Unmarshal(jsonBody, &module)
	if err != nil {
		fmt.Println("error parsing Get response JSON:", err)
	}

	return module, err
}

func (smilecdr *Client) PostModuleConfig(ctx context.Context, nodeId string, module ModuleConfig) (ModuleConfig, error) {
	var newModule ModuleConfig
	var moduleId = module.ModuleId

	var endpoint = fmt.Sprintf("/module-config/%s/%s/create", nodeId, moduleId)
	jsonBody, _ := json.Marshal(module)

	fmt.Println("PostModuleConfig: ", string(jsonBody))

	jsonBody, postErr := smilecdr.Post(ctx, endpoint, jsonBody)
	if postErr != nil {
		fmt.Println("error during Post in PostOpenIdClient:", postErr)
		return newModule, postErr
	}

	err := json.Unmarshal(jsonBody, &newModule)
	if err != nil {
		fmt.Println("error parsing Post response JSON:", err)
	}

	return newModule, err
}

func (smilecdr *Client) PutModuleConfig(ctx context.Context, nodeId string, module ModuleConfig) (ModuleConfig, error) {
	var moduleId = module.ModuleId

	var endpoint = fmt.Sprintf("/module-config/%s/%s/set", nodeId, moduleId)

	printModuleConfig(module) // DEBUG

	jsonBody, error := json.Marshal(module)
	if error != nil {
		fmt.Println("error marshalling moduleConfig:", error)
		return module, error
	}

	fmt.Println("PutModuleConfig: ", string(jsonBody))

	resp, putErr := smilecdr.Put(ctx, endpoint, jsonBody)

	if putErr != nil {
		fmt.Println("error during Put in PutModuleConfig:", putErr)
		fmt.Println("ResponseBody:", string(resp))
		return module, putErr
	}

	return module, nil
}

func (smilecdr *Client) DeleteModuleConfig(ctx context.Context, nodeId string, moduleId string) error {
	var endpoint = fmt.Sprintf("/module-config/%s/%s/archive", nodeId, moduleId)
	_, err := smilecdr.Delete(ctx, endpoint)
	return err
}
