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

// For Integrations tests, the ActionTests struct is referenced from testdata.go which needs to be created based on testdata.go.template.
// Integrations tests will only run if the environment variable `INTEGRATION` is set.

package turboclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientAuth_BasicAuth(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/v3/login", r.URL.Path)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, "username=testuser&password=testpass", string(body))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"uuid":"1234567890","username":"administrator",
		"roles":[{"name":"ADMINISTRATOR"}],"loginProvider":"Local",
		"authToken":"","showSharedUserSC":false}`))
		}))
	defer server.Close()

	authReq := AuthRequest{
		basePath:   "/api/v3",
		hostname:   strings.Replace(server.URL, "https://", "", 1),
		username:   "testuser",
		password:   "testpass",
		httpClient: server.Client(),
	}

	client, err := clientAuth(&authReq)
	assert.NoError(t, err)
	assert.Equal(t, server.URL+"/api/v3", client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
}

func TestClientAuth_OAuth2(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/oauth2/token", r.URL.Path)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, "grant_type=client_credentials&scope=role:OBSERVER", string(body))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token": "admin_token",
			"scope": "read write", "token_type": "Bearer", "expires_in": 3600}`))
		}))
	defer server.Close()

	authReq := AuthRequest{
		basePath: "/api/v3",
		hostname: strings.Replace(server.URL, "https://", "", 1),
		oAuthCreds: OAuthCreds{
			ClientId:     "test_client",
			ClientSecret: "test_secret",
			Role:         OBSERVER,
		},
		httpClient: server.Client(),
	}

	client, err := clientAuth(&authReq)
	assert.NoError(t, err)
	assert.Equal(t, server.URL+"/api/v3", client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, "Bearer admin_token", client.Headers["Authorization"])
}

func TestClientAuth_InvalidCredentials(t *testing.T) {
	authReq := AuthRequest{
		basePath: "/api/v3",
		hostname: "127.0.0.1:12345",
		httpClient: httptest.NewTLSServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			})).Client(),
	}

	_, err := clientAuth(&authReq)
	assert.Error(t, err)
}
