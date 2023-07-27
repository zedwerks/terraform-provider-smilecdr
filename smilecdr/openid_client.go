// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"encoding/json"
	"fmt"
)

type ClientSecret struct {
	Pid         int    `json:"pid,omitempty"`
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
	AccessTokenValiditySeconds  int              `json:"accessTokenValiditySeconds,omitempty"`
	AllowedGrantTypes           []string         `json:"allowedGrantTypes,omitempty"`
	AutoApproveScopes           []string         `json:"autoApproveScopes,omitempty"`
	AutoGrantScopes             []string         `json:"autoGrantScopes,omitempty"`
	ClientId                    string           `json:"clientId,omitempty"`
	ClientName                  string           `json:"clientName,omitempty"`
	ClientSecrets               []ClientSecret   `json:"clientSecrets,omitempty"`
	FixedScope                  bool             `json:"fixedScope,omitempty"`
	RefreshTokenValiditySeconds int              `json:"refreshTokenValiditySeconds,omitempty"`
	RegisteredRedirectUris      []string         `json:"registeredRedirectUris,omitempty"`
	Scopes                      []string         `json:"scopes,omitempty"`
	SecretRequired              bool             `json:"secretRequired,omitempty"`
	SecretClientCanChange       bool             `json:"secretClientCanChange,omitempty"`
	Enabled                     bool             `json:"enabled,omitempty"`
	CanIntrospectAnyTokens      bool             `json:"canIntrospectAnyTokens,omitempty"`
	CanIntrospectOwnTokens      bool             `json:"canIntrospectOwnTokens,omitempty"`
	AlwaysRequireApproval       bool             `json:"alwaysRequireApproval,omitempty"`
	CanReissueTokens            bool             `json:"canReissueTokens,omitempty"`
	Permissions                 []UserPermission `json:"permissions,omitempty"`
	AttestationAccepted         bool             `json:"rememberedScopes,omitempty"`
	PublicJwksUri               string           `json:"publicJwksUri,omitempty"`
	ArchivedAt                  string           `json:"archivedAt,omitempty"`
	CreatedByAppSphere          bool             `json:"createdByAppSphere,omitempty"`
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
