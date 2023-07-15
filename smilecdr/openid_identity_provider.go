// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"encoding/json"
	"fmt"
)

type OpenIdIdentityProvider struct {
	Pid                             int    `json:"pid,omitempty"`
	Name                            string `json:"name,omitempty"`
	Issuer                          string `json:"issuer,omitempty"`
	TokenIntrospectionClientId      string `json:"tokenIntrospectionClientId,omitempty"`
	TokenIntrospectionClientSecret  string `json:"tokenIntrospectionClientSecret,omitempty"`
	NodeId                          string `json:"node_id,omitempty"`
	ModuleId                        string `json:"module_id,omitempty"`
	ValidationJwkText               string `json:"validation_jwk_text,omitempty"`
	ValidationJwkFile               string `json:"validation_jwk_file,omitempty"`
	FederationRegistrationId        string `json:"federation_registration_id,omitempty"`
	FederationRequestScopes         string `json:"federation_request_scopes,omitempty"`
	FederationAuthorizationUrl      string `json:"federation_authorization_url,omitempty"`
	FederationTokenUrl              string `json:"federation_token_url,omitempty"`
	FederationUserInfoUrl           string `json:"federation_user_info_url,omitempty"`
	FederationJwkSetUrl             string `json:"federation_jwk_set_url,omitempty"`
	FederationAuthScriptText        string `json:"federation_auth_script_text,omitempty"`
	FederationAuthScriptFile        string `json:"federation_auth_script_file,omitempty"`
	FederationUserMappingScriptText string `json:"federation_user_mapping_script_text,omitempty"`
	FhirEndpointUrl                 string `json:"fhir_endpoint_url,omitempty"`
	AuthWellKnownConfigUrl          string `json:"auth_well_known_config_url,omitempty"`
	Notes                           string `json:"notes,omitempty"`
	CustomTokenParams               string `json:"custom_token_params,omitempty"`
	ResponseType                    string `json:"response_type,omitempty"`
	OrganizationId                  string `json:"organization_id,omitempty"`
	Audience                        string `json:"audience,omitempty"`
	ArchivedAt                      string `json:"archived_at,omitempty"`
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

func (smilecdr *Client) GetOpenIdIdentityProvider(nodeId string, moduleId string, name string) (OpenIdIdentityProvider, error) {

	var provider OpenIdIdentityProvider
	var endpoint = fmt.Sprintf("/openid-connect-servers/%s/%s/%s", nodeId, moduleId, name)
	jsonBody, getErr := smilecdr.Get(endpoint)
	if getErr != nil {
		fmt.Println("error during Get in GetOpenIdIdentityProvider:", getErr)
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
		fmt.Println("error during Post in PostOpenIdIdentityProvider:", postErr)
		return newProvider, postErr
	}

	err := json.Unmarshal(jsonBody, &newProvider)

	return newProvider, err
}
