// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"encoding/json"
	"fmt"
)

type ClientSecret struct {
	Secret      string `json:"secret,omitempty"`
	Description string `json:"description,omitempty"`
	Expiration  string `json:"expiration,omitempty"`
	Activation  string `json:"activation,omitempty"`
}

type UserPermission struct {
	Permission string `json:"permission,omitempty"`
	Argument   string `json:"argument,omitempty"`
}

type OpenIdClient struct {
	Pid                         int              `json:"pid,omitempty"`
	NodeId                      string           `json:"nodeId,omitempty"`
	ModuleId                    string           `json:"moduleId,omitempty"`
	ClientId                    string           `json:"clientId,omitempty"`
	ClientName                  string           `json:"clientName,omitempty"`
	Enabled                     bool             `json:"enabled,omitempty"`
	AccessTokenValiditySeconds  int              `json:"accessTokenValiditySeconds,omitempty"`
	AllowedGrantTypes           []string         `json:"allowedGrantTypes,omitempty"`
	AlwaysRequireApproval       bool             `json:"alwaysRequireApproval,omitempty"`
	AttestationAccepted         bool             `json:"attestationAccepted,omitempty"`
	AutoApproveScopes           []string         `json:"autoApproveScopes,omitempty"`
	AutoGrantScopes             []string         `json:"autoGrantScopes,omitempty"`
	CanIntrospectAnyTokens      bool             `json:"canIntrospectAnyTokens,omitempty"`
	CanIntrospectOwnTokens      bool             `json:"canIntrospectOwnTokens,omitempty"`
	CanReissueTokens            bool             `json:"canReissueTokens,omitempty"`
	ClientSecrets               []ClientSecret   `json:"clientSecrets,omitempty"`
	CreatedByAppSphere          bool             `json:"createdByAppSphere,omitempty"`
	FixedScope                  bool             `json:"fixedScope,omitempty"`
	JwksUrl                     string           `json:"jwksUrl,omitempty"`
	Permissions                 []UserPermission `json:"permissions,omitempty"`
	PublicJwks                  string           `json:"publicJwks,omitempty"`
	RefreshTokenValiditySeconds int              `json:"refreshTokenValiditySeconds,omitempty"`
	RegisteredRedirectUris      []string         `json:"registeredRedirectUris,omitempty"`
	RememberApprovedScopes      bool             `json:"rememberApprovedScopes,omitempty"`
	Scopes                      []string         `json:"scopes,omitempty"`
	SecretClientCanChange       bool             `json:"secretClientCanChange,omitempty"`
	SecretRequired              bool             `json:"secretRequired,omitempty"`
	ArchivedAt                  string           `json:"archivedAt,omitempty"`
}

func (smilecdr *Client) GetOpenIdClients() ([]OpenIdClient, error) {
	var clients []OpenIdClient
	jsonBody, getErr := smilecdr.Get("/openid-connect-clients")
	if getErr != nil {
		return clients, getErr
	}

	err := json.Unmarshal(jsonBody, &clients)

	return clients, err
}

func (smilecdr *Client) GetOpenIdClient(nodeId string, moduleId string, clientId string) (OpenIdClient, error) {
	var client OpenIdClient
	var endpoint = fmt.Sprintf("/openid-connect-clients/%s/%s/%s", nodeId, moduleId, clientId)
	jsonBody, getErr := smilecdr.Get(endpoint)
	if getErr != nil {
		fmt.Println("error during Get in GetOpenIdClient:", getErr)
		return client, getErr
	}

	err := json.Unmarshal(jsonBody, &client)
	if err != nil {
		fmt.Println("error parsing Get response JSON:", err)
	}

	return client, err
}

func (smilecdr *Client) PostOpenIdClient(client OpenIdClient) (OpenIdClient, error) {
	var newClient OpenIdClient
	var nodeId = client.NodeId
	var moduleId = client.ModuleId

	var endpoint = fmt.Sprintf("/openid-connect-clients/%s/%s", nodeId, moduleId)
	jsonBody, _ := json.Marshal(client)

	fmt.Println("jsonBody:", string(jsonBody))

	jsonBody, postErr := smilecdr.Post(endpoint, jsonBody)
	if postErr != nil {
		fmt.Println("error during Post in PostOpenIdClient:", postErr)
		return newClient, postErr
	}

	err := json.Unmarshal(jsonBody, &newClient)
	if err != nil {
		fmt.Println("error parsing Post response JSON:", err)
	}

	return newClient, err
}

func (smilecdr *Client) PutOpenIdClient(client OpenIdClient) (OpenIdClient, error) {
	var newClient OpenIdClient
	var nodeId = client.NodeId
	var moduleId = client.ModuleId
	var clientId = client.ClientId

	var endpoint = fmt.Sprintf("/openid-connect-clients/%s/%s/%s", nodeId, moduleId, clientId)

	jsonBody, _ := json.Marshal(client)

	fmt.Println("jsonBody:", string(jsonBody))

	jsonBody, putErr := smilecdr.Put(endpoint, jsonBody)
	if putErr != nil {
		fmt.Println("error during Put in PutOpenIdClient:", putErr)
		return newClient, putErr
	}

	err := json.Unmarshal(jsonBody, &newClient)
	if err != nil {
		fmt.Println("error parsing Put response JSON:", err)
	}

	return newClient, err
}

func (smilecdr *Client) DeleteOpenIdClient(nodeId string, moduleId string, clientId string) error {
	var endpoint = fmt.Sprintf("/openid-connect-clients/%s/%s/%s", nodeId, moduleId, clientId)
	_, err := smilecdr.Delete(endpoint)

	return err
}
