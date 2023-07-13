// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"bytes"
	"encoding/base64"
	"fmt"
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
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error making HTTP Get Request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("received non-200 OK status code:", resp.StatusCode)
		// Handle the error condition here
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading HTTP Get Response Body:", err)
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
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error making HTTP Post Request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("received non-200 OK status code:", resp.StatusCode)
		// Handle the error condition here
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading HTTP Post Response Body:", err)
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
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("received non-200 OK status code:", resp.StatusCode)
		// Handle the error condition here
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading HTTP Put Response Body:", err)
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
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error making HTTP Delete Request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("received non-200 OK status code:", resp.StatusCode)
		// Handle the error condition here
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading HTTP Delete Response Body:", err)
		return nil, err
	}

	return body, nil
}
