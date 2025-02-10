// Copyright (c) IBM Corporation
// SPDX-License-Identifier: Apache-2.0

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS-IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package turboclient

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

// Parameters for authenticating to the Turbonomic API
type AuthRequest struct {
	basePath   string
	hostname   string
	username   string
	password   string
	httpClient *http.Client
}

// Creates authorized Turbonomic API Client
func clientAuth(authreq *AuthRequest) (*Client, error) {
	var client Client

	urlPath := "https://" + authreq.hostname + authreq.basePath + "/login"
	payload := strings.NewReader("username=" + authreq.username + "&password=" + authreq.password)

	req, err := http.NewRequest("POST", urlPath, payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := authreq.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		slog.Debug("Failed to establish a connection with the Turbonomic instance:", "Status", string(resp.Status))
		err := errors.New(resp.Status)
		return nil, err
	}
	slog.Debug("Successfully logged into Turbonomic")

	client.BaseURL = "https://" + authreq.hostname + authreq.basePath
	client.HTTPClient = authreq.httpClient

	return &client, nil
}
