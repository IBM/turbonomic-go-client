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
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
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
	oAuthCreds OAuthCreds
	httpClient *http.Client
	apiInfo    ApiInfo
}

type oAuthResp struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Creates authorized Turbonomic API Client
func clientAuth(authreq *AuthRequest) (*Client, error) {
	var client Client
	if authreq == nil {
		return nil, errors.New("please provide valid credentials")
	}
	urlPath, payload, err := setAuthParams(*authreq)

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", urlPath, payload)

	if err != nil {
		return nil, err
	}

	var apiOrigin = authreq.apiInfo.ApiOrigin
	var version = authreq.apiInfo.Version
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if apiOrigin != "" {
		req.Header.Add("User-Agent", apiOrigin+"/"+version)
	}
	if (authreq.oAuthCreds.ClientId != "") && (authreq.oAuthCreds.ClientSecret != "") {
		auth := authreq.oAuthCreds.ClientId + ":" + authreq.oAuthCreds.ClientSecret
		req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}

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

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var result oAuthResp
	if err := json.Unmarshal(body, &result); err != nil {
		slog.Error("Can not unmarshal oAuth JSON Response")
		return nil, err
	}

	client.BaseURL = "https://" + authreq.hostname + authreq.basePath
	client.HTTPClient = authreq.httpClient
	client.Headers = make(map[string]string)
	if result.AccessToken != "" {
		client.Headers["Authorization"] = "Bearer " + result.AccessToken
	}
	if apiOrigin != "" {
		client.Headers["User-Agent"] = apiOrigin + "/" + version
	}
	return &client, nil
}

func setAuthParams(authreq AuthRequest) (urlPath string, payload *strings.Reader, err error) {
	if (authreq.username) != "" && (authreq.password != "") {
		urlPath = "https://" + authreq.hostname + authreq.basePath + "/login"
		payload = strings.NewReader("username=" + authreq.username + "&password=" + authreq.password)
	} else if (authreq.oAuthCreds.Role.String() != "") &&
		(authreq.oAuthCreds.ClientId != "") &&
		(authreq.oAuthCreds.ClientSecret != "") {

		urlPath = "https://" + authreq.hostname + "/oauth2/token"
		payload = strings.NewReader("grant_type=" + "client_credentials" +
			"&scope=role:" + authreq.oAuthCreds.Role.String())
	} else {
		return "", nil, errors.New("please provide valid credentials; username/password or oauth2")
	}
	return urlPath, payload, nil

}
