// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
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
	uri := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Accept", "application/json")

	fmt.Println("GET Request URI: ", uri)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error making HTTP Get Request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Http GET: received non-200 OK status code:", resp.StatusCode)
		// Handle the error condition here
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Http GET: error reading Response Body:", err)
		return nil, err
	}

	return body, nil
}

func (c *Client) Post(endpoint string, body []byte) ([]byte, error) {
	uri := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("Http Post: error during request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		fmt.Println("Post(): Did not receive 200, 201 or 204. Received this status code:", resp.StatusCode)
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading HTTP Post Response Body:", err)
		return nil, err
	}

	return body, nil
}

func (c *Client) Put(endpoint string, body []byte) ([]byte, error) {
	uri := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	fmt.Println("PUT Request URI: ", uri)
	fmt.Println("PUT Request Body: ", string(body))

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

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading HTTP Put Response Body:", err)
		return nil, err
	}

	return body, nil
}

func (c *Client) Delete(endpoint string) ([]byte, error) {
	uri := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Accept", "application/json")

	fmt.Println("DELETE Request URI: ", uri)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("error making HTTP Delete Request: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		fmt.Println("Expected 200, or 204, instead received status code:", resp.StatusCode)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading HTTP Delete Response Body:", err)
		return nil, err
	}

	return body, nil
}
