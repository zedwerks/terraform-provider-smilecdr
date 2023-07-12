// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseUrl    string
	authHeader string
	httpClient *http.Client
}

func NewClient(baseUrl string, username string, password string) *Client {
	credentials := username + ":" + password
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))

	return &Client{
		baseUrl:    baseUrl,
		authHeader: auth,
		httpClient: &http.Client{},
	}
}

func (c *Client) Get(endpoint string) ([]byte, error) {
	url := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) Post(endpoint string, body []byte) ([]byte, error) {
	url := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) Put(endpoint string, body []byte) ([]byte, error) {
	url := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) Delete(endpoint string) ([]byte, error) {
	url := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
