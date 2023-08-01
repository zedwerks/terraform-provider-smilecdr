// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"encoding/json"
	"fmt"
)

type ModuleConfig struct {
	ModuleId     string               `json:"moduleId,omitempty"`
	ModuleType   string               `json:"moduleType,omitempty"`
	Options      []ModuleOption       `json:"options,omitempty"`
	Dependencies []ModuleDependencies `json:"dependencies,omitempty"`
}

type ModuleOption struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type ModuleDependencies struct {
	ModuleId string `json:"moduleId,omitempty"`
	Type     string `json:"type,omitempty"`
}

func (moduleConfig *ModuleConfig) LookupOption(key string) (string, error) {

	moduleConfigOptions := moduleConfig.Options
	for _, kv := range moduleConfigOptions {
		if kv.Key == key {
			return kv.Value, nil
		}
	}

	return "", fmt.Errorf("value not found for key '%s'", key)
}

func (smilecdr *Client) GetModuleConfigs() ([]ModuleConfig, error) {
	var modules []ModuleConfig
	jsonBody, getErr := smilecdr.Get("/module-config")
	if getErr != nil {
		return modules, getErr
	}

	err := json.Unmarshal(jsonBody, &modules)

	return modules, err
}

func (smilecdr *Client) GetModuleConfig(moduleId string) (ModuleConfig, error) {
	var module ModuleConfig
	var endpoint = fmt.Sprintf("/module-config/%s", moduleId)
	jsonBody, getErr := smilecdr.Get(endpoint)
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

func (smilecdr *Client) PostModuleConfig(module ModuleConfig) (ModuleConfig, error) {
	var newModule ModuleConfig
	var moduleId = module.ModuleId

	var endpoint = fmt.Sprintf("/module-config/%s", moduleId)
	jsonBody, _ := json.Marshal(module)

	jsonBody, postErr := smilecdr.Post(endpoint, jsonBody)
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

func (smilecdr *Client) PutModuleConfig(module ModuleConfig) (ModuleConfig, error) {
	var newModule ModuleConfig
	var moduleId = module.ModuleId

	var endpoint = fmt.Sprintf("/module-config/%s", moduleId)
	jsonBody, _ := json.Marshal(module)

	jsonBody, putErr := smilecdr.Put(endpoint, jsonBody)
	if putErr != nil {
		fmt.Println("error during Put in PutModuleConfig:", putErr)
		return newModule, putErr
	}

	err := json.Unmarshal(jsonBody, &newModule)
	if err != nil {
		fmt.Println("error parsing Put response JSON:", err)
	}

	return newModule, err
}

func (smilecdr *Client) DeleteModuleConfig(moduleId string) error {
	var endpoint = fmt.Sprintf("/module-config/%s", moduleId)
	_, err := smilecdr.Delete(endpoint)
	return err
}
