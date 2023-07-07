// Copyright (c) Zed Werks Inc.
// SPDX-License-Identifier: APACHE-2.0

package smilecdr

import (
	"encoding/base64"
	"net/http"
)

type authTransport struct {
	Transport http.RoundTripper
	header    string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.header)
	return t.Transport.RoundTrip(req)
}

func newClient(username, password string) *http.Client {
	credentials := username + ":" + password
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	// Add the Authorization header to every request
	client.Transport = &authTransport{
		Transport: client.Transport,
		header:    authHeader,
	}

	return client
}
