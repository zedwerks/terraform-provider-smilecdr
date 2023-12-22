// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Client struct {
	baseUrl    string
	authHeader string
	httpClient *http.Client
	debug      bool
}

func NewClient(ctx context.Context, baseUrl string, username string, password string) *Client {
	credentials := username + ":" + password
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))

	if baseUrl == "" {
		tflog.Warn(ctx, "Missing SmileCDR Admin API Base Url.")
	}
	tfLog, _ := os.LookupEnv("TF_LOG")

	smilecdrClient := Client{
		baseUrl:    baseUrl,
		authHeader: auth,
		httpClient: &http.Client{},
		debug:      (tfLog == "DEBUG"),
	}
	return &smilecdrClient
}

func (c *Client) Get(ctx context.Context, endpoint string) ([]byte, error) {
	uri := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Accept", "application/json")

	logArgs := map[string]interface{}{
		"uri": uri,
	}

	tflog.Debug(ctx, "Http GET Request URI:\n", logArgs)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Http GET: error during request: %s", err.Error())
		tflog.Error(ctx, errMsg)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Http GET: received non-200 OK status code: %d", resp.StatusCode)
		tflog.Error(ctx, errMsg)
		return nil, err
	}

	rBody, err := io.ReadAll(resp.Body)

	if err != nil {
		errMsg := fmt.Sprintf("Http GET: error reading Response Body: %s", err.Error())
		tflog.Error(ctx, errMsg)
		return nil, err
	}

	return rBody, nil
}

func (c *Client) Post(ctx context.Context, endpoint string, body []byte) ([]byte, error) {
	uri := c.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		errMsg := fmt.Sprintf("Http POST: error creating request: %s", err.Error())
		tflog.Error(ctx, errMsg)
		return nil, err
	}
	req.Header.Add("Authorization", c.authHeader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Http POST: error during request: %s", err.Error())
		tflog.Error(ctx, errMsg)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		errMsg := fmt.Sprintf("Http POST: Expecting 200, 201 or 204. Received: %d", resp.StatusCode)
		tflog.Info(ctx, errMsg)
		return nil, err
	}

	var rBody []byte = nil

	rBody, err = io.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Http POST: Error reading Response Body: %s", err)
		tflog.Info(ctx, errMsg)
		return nil, err
	}

	return rBody, nil
}

func (c *Client) Put(ctx context.Context, endpoint string, body []byte) ([]byte, error) {
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
		fmt.Printf("[ERROR] Http PUT: error during request: %s", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[WARN] Http PUT: Received non-200 status code: %d", resp.StatusCode)
		// Handle the error condition here
		return nil, err
	}

	var rBody []byte = nil

	rBody, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[ERROR] Http PUT: Error reading Response Body: %s", err.Error())
		return nil, err
	}

	return rBody, nil
}

func (c *Client) Delete(ctx context.Context, endpoint string) ([]byte, error) {
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
		fmt.Printf("[ERROR] Http DELETE Error during Request: %s", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		fmt.Printf("[WARN] Http DELETE: Expected 200, or 204, instead received status code: %d", resp.StatusCode)
		return nil, err
	}

	rBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[ERROR] Http DELETE: Error reading Response Body:", err)
		return nil, err
	}

	return rBody, nil
}
