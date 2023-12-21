// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"context"
	"encoding/json"
	"fmt"
)

type UserAuthorities struct {
	Permission string `json:"permission,omitempty"`
	Argument   string `json:"argument,omitempty"`
}

type User struct {
	Pid                 int               `json:"pid,omitempty"`
	NodeId              string            `json:"nodeId,omitempty"`
	ModuleId            string            `json:"moduleId,omitempty"`
	Username            string            `json:"username,omitempty"`
	Password            string            `json:"password,omitempty"`
	FamilyName          string            `json:"familyName,omitempty"`
	GivenName           string            `json:"givenName,omitempty"`
	AccountLocked       bool              `json:"accountLocked,omitempty"`
	SystemUser          bool              `json:"systemUser,omitempty"`
	AccountDisabled     bool              `json:"accountDisabled,omitempty"`
	External            bool              `json:"external,omitempty"`
	ServiceAccount      bool              `json:"serviceAccount,omitempty"`
	TwoFactorAuthStatus string            `json:"twoFactorAuthStatus,omitempty"`
	Authorities         []UserAuthorities `json:"authorities,omitempty"`

	// These fields are Computed
	LastConnected string `json:"lastConnected,omitempty"`
	LastActive    string `json:"lastActive,omitempty"`
}

func (smilecdr *Client) GetUser(ctx context.Context, nodeId string, moduleId string, pid int) (User, error) {
	var user User
	var err error
	var endpoint = fmt.Sprintf("/user-management/%s/%s/%d", nodeId, moduleId, pid)
	jsonBody, err := smilecdr.Get(ctx, endpoint)
	if err != nil {
		fmt.Println("error during Get in GetUser:", err)
		return user, err
	}

	if jsonBody != nil {
		fmt.Println("GET Response jsonBody:", string(jsonBody))
		err = json.Unmarshal(jsonBody, &user)
		if err != nil {
			fmt.Println("error parsing Get response JSON:", err)
		}
	}
	return user, err
}

func (smilecdr *Client) PostUser(ctx context.Context, user User) (User, error) {
	var newUser User
	var nodeId = user.NodeId
	var moduleId = user.ModuleId

	var endpoint = fmt.Sprintf("/user-management/%s/%s", nodeId, moduleId)
	jsonBody, _ := json.Marshal(user)

	fmt.Println("POST Request jsonBody:", string(jsonBody))

	jsonBody, postErr := smilecdr.Post(ctx, endpoint, jsonBody)
	if postErr != nil {
		fmt.Println("error during Post in PostUser:", postErr)
		return newUser, postErr
	}

	err := json.Unmarshal(jsonBody, &newUser)
	if err != nil {
		fmt.Println("error parsing Post response JSON:", err)
	}

	return newUser, err
}

func (smilecdr *Client) PutUser(ctx context.Context, user User) (User, error) {
	var updatedUser User
	var nodeId = user.NodeId
	var moduleId = user.ModuleId

	var endpoint = fmt.Sprintf("/user-management/%s/%s/%d", nodeId, moduleId, user.Pid)
	jsonBody, _ := json.Marshal(user)

	jsonBody, putErr := smilecdr.Put(ctx, endpoint, jsonBody)
	if putErr != nil {
		fmt.Println("error during Put in PutUser:", putErr)
		return updatedUser, putErr
	}

	err := json.Unmarshal(jsonBody, &updatedUser)
	if err != nil {
		fmt.Println("error parsing Put response JSON:", err)
	}

	return updatedUser, err
}

func (smilecdr *Client) DeleteUser(ctx context.Context, nodeId string, moduleId string, pid int) error {
	var endpoint = fmt.Sprintf("/user-management/%s/%s/%d", nodeId, moduleId, pid)
	_, err := smilecdr.Delete(ctx, endpoint)
	if err != nil {
		fmt.Println("error during Delete in DeleteUser:", err)
	}
	return err
}
