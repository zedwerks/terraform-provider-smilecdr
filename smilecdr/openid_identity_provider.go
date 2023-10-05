// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type OpenIdIdentityProvider struct {
	Pid                             int    `json:"pid,omitempty"`
	Name                            string `json:"name,omitempty"`
	Issuer                          string `json:"issuer,omitempty"`
	TokenIntrospectionClientId      string `json:"tokenIntrospectionClientId,omitempty"`
	TokenIntrospectionClientSecret  string `json:"tokenIntrospectionClientSecret,omitempty"`
	NodeId                          string `json:"nodeId,omitempty"`
	ModuleId                        string `json:"moduleId,omitempty"`
	ValidationJwkText               string `json:"validationJwkText,omitempty"`
	ValidationJwkFile               string `json:"validationJwkFile,omitempty"`
	FederationRegistrationId        string `json:"federationRegistrationId,omitempty"`
	FederationRequestScopes         string `json:"federationRequestScopes,omitempty"`
	FederationAuthorizationUrl      string `json:"federationAuthorizationUrl,omitempty"`
	FederationTokenUrl              string `json:"federationTokenUrl,omitempty"`
	FederationUserInfoUrl           string `json:"federationUserInfoUrl,omitempty"`
	FederationJwkSetUrl             string `json:"federationJwkSetUrl,omitempty"`
	FederationAuthScriptText        string `json:"federationAuthScriptText,omitempty"`
	FederationUserMappingScriptText string `json:"federationUserMappingScriptText,omitempty"`
	ArchivedAt                      string `json:"archivedAt,omitempty"`
}

func (smilecdr *Client) GetOpenIdIdentityProviders() ([]OpenIdIdentityProvider, error) {
	var providers []OpenIdIdentityProvider
	jsonBody, getErr := smilecdr.Get("/openid-connect-servers")
	if getErr != nil {
		return providers, getErr
	}

	err := json.Unmarshal(jsonBody, &providers)

	return providers, err
}

func (smilecdr *Client) GetOpenIdIdentityProvider(nodeId string, moduleId string, issuerUrl string) (OpenIdIdentityProvider, error) {

	var provider OpenIdIdentityProvider
	var endpoint = fmt.Sprintf("/openid-connect-servers/%s/%s/?issuer_url=%s", nodeId, moduleId, url.PathEscape(issuerUrl))
	jsonBody, getErr := smilecdr.Get(endpoint)
	if getErr != nil {
		fmt.Println("error during GET in GetOpenIdIdentityProvider:", getErr)
		return provider, getErr
	}

	err := json.Unmarshal(jsonBody, &provider)

	return provider, err
}

func (smilecdr *Client) PostOpenIdIdentityProvider(provider OpenIdIdentityProvider) (OpenIdIdentityProvider, error) {
	var newProvider OpenIdIdentityProvider
	var nodeId = provider.NodeId
	var moduleId = provider.ModuleId

	var endpoint = fmt.Sprintf("/openid-connect-servers/%s/%s", nodeId, moduleId)
	jsonBody, _ := json.Marshal(provider)

	jsonBody, postErr := smilecdr.Post(endpoint, jsonBody)
	if postErr != nil {
		fmt.Println("error during POST in PostOpenIdIdentityProvider:", postErr)
		return newProvider, postErr
	}

	err := json.Unmarshal(jsonBody, &newProvider)

	return newProvider, err
}

func (smilecdr *Client) PutOpenIdIdentityProvider(provider OpenIdIdentityProvider) (OpenIdIdentityProvider, error) {
	var nodeId = provider.NodeId
	var moduleId = provider.ModuleId
	var pid = provider.Pid

	var endpoint = fmt.Sprintf("/openid-connect-servers/%s/%s/%s", nodeId, moduleId, strconv.Itoa(pid))
	jsonBody, _ := json.Marshal(provider)

	resp, putErr := smilecdr.Put(endpoint, jsonBody)
	if putErr != nil {
		fmt.Println("error during PUT in PutOpenIdIdentityProvider:", putErr)
		fmt.Println("ResponseBody:", string(resp))
		return provider, putErr
	}

	return provider, nil
}
