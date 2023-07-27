// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import "encoding/json"

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

func (smilecdr *Client) GetModuleConfigs() ([]ModuleConfig, error) {
	var modules []ModuleConfig
	jsonBody, getErr := smilecdr.Get("/module-config")
	if getErr != nil {
		return modules, getErr
	}

	err := json.Unmarshal(jsonBody, &modules)

	return modules, err
}
