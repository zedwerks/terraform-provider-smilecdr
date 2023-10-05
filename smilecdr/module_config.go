// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
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

func (moduleConfig *ModuleConfig) LookupOptionOk(key string) (string, bool) {

	moduleConfigOptions := moduleConfig.Options
	for _, kv := range moduleConfigOptions {
		if kv.Key == key {
			return kv.Value, true
		}
	}

	return "", false
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

func (smilecdr *Client) GetModuleConfig(nodeId string, moduleId string) (ModuleConfig, error) {
	var module ModuleConfig
	var endpoint = fmt.Sprintf("/module-config/%s/%s", nodeId, moduleId)
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

func (smilecdr *Client) PostModuleConfig(nodeId string, module ModuleConfig) (ModuleConfig, error) {
	var newModule ModuleConfig
	var moduleId = module.ModuleId

	var endpoint = fmt.Sprintf("/module-config/%s/%s/create", nodeId, moduleId)
	jsonBody, _ := json.Marshal(module)

	fmt.Println("PostModuleConfig: ", string(jsonBody))

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

func (smilecdr *Client) PutModuleConfig(nodeId string, module ModuleConfig) (ModuleConfig, error) {
	var moduleId = module.ModuleId

	var endpoint = fmt.Sprintf("/module-config/%s/%s/set", nodeId, moduleId)
	jsonBody, _ := json.Marshal(module)

	resp, putErr := smilecdr.Put(endpoint, jsonBody)
	if putErr != nil {
		fmt.Println("error during Put in PutModuleConfig:", putErr)
		fmt.Println("ResponseBody:", string(resp))
		return module, putErr
	}

	return module, nil
}

func (smilecdr *Client) DeleteModuleConfig(nodeId string, moduleId string) error {
	var endpoint = fmt.Sprintf("/module-config/%s/%s/archive", nodeId, moduleId)
	_, err := smilecdr.Delete(endpoint)
	return err
}
