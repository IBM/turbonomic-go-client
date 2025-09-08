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
	"github.com/IBM/turbonomic-go-client/logging"
)

func TestClientAuth_BasicAuth(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/v3/login", r.URL.Path)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, "username=testuser&password=testpass", string(body))
			assert.Equal(t, "Go-http-client/1.1", r.Header.Get("User-Agent"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"uuid":"1234567890","username":"administrator",
		"roles":[{"name":"ADMINISTRATOR"}],"loginProvider":"Local",
		"authToken":"","showSharedUserSC":false}`)); err != nil {
				t.Fail()
				t.Log(err)
			}
		}))
	defer server.Close()

	authReq := AuthRequest{
		basePath:   "/api/v3",
		hostname:   strings.Replace(server.URL, "https://", "", 1),
		username:   "testuser",
		password:   "testpass",
		httpClient: server.Client(),
	}

	emptyLogConfig := logging.SetLogConfig([]logging.LoggingOption{})
	client, err := clientAuth(&authReq, emptyLogConfig)
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
			if _, err := w.Write([]byte(`{"access_token": "admin_token",
			"scope": "read write", "token_type": "Bearer", "expires_in": 3600}`)); err != nil {
				t.Fail()
				t.Log(err)
			}
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

	emptyLogConfig := logging.SetLogConfig([]logging.LoggingOption{})
	client, err := clientAuth(&authReq, emptyLogConfig)
	assert.NoError(t, err)
	assert.Equal(t, server.URL+"/api/v3", client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, "Bearer admin_token", client.Headers["Authorization"])
}

func TestClientAuth_OAuth2POST(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "error reading body", http.StatusBadRequest)
			}
			defer r.Body.Close()
			bodyStr := string(body)

			switch {
			case r.URL.Path == "/oauth2/token" && r.Method == "POST" && !strings.Contains(bodyStr, "client_id"):
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "/oauth2/token", r.URL.Path)
				assert.Equal(t, "grant_type=client_credentials&scope=role:OBSERVER", string(body))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
			default:
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "/oauth2/token", r.URL.Path)
				assert.Equal(t, "grant_type=client_credentials&scope=role:OBSERVER&client_id=test_client&client_secret=test_secret", string(body))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(`{"access_token": "admin_token",
			"scope": "read write", "token_type": "Bearer", "expires_in": 3600}`)); err != nil {
					t.Fail()
					t.Log(err)
				}
			}
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

	emptyLogConfig := logging.SetLogConfig([]logging.LoggingOption{})
	client, err := clientAuth(&authReq, emptyLogConfig)
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
	emptyLogConfig := logging.SetLogConfig([]logging.LoggingOption{})
	_, err := clientAuth(&authReq, emptyLogConfig)
	assert.Error(t, err)
}

func TestClientAuth_ProviderApiInfo(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/api/v3/login", r.URL.Path)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, "username=testuser&password=testpass", string(body))
			assert.Equal(t, "turbonomic-terraform-provider/1.1.0", r.Header.Get("User-Agent"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"uuid":"1234567890","username":"administrator",
		"roles":[{"name":"ADMINISTRATOR"}],"loginProvider":"Local",
		"authToken":"","showSharedUserSC":false}`)); err != nil {
				t.Fail()
				t.Log(err)
			}
		}))
	defer server.Close()

	authReq := AuthRequest{
		basePath:   "/api/v3",
		hostname:   strings.Replace(server.URL, "https://", "", 1),
		username:   "testuser",
		password:   "testpass",
		httpClient: server.Client(),
		apiInfo: ApiInfo{
			ApiOrigin: "turbonomic-terraform-provider",
			Version:   "1.1.0",
		},
	}
	emptyLogConfig := logging.SetLogConfig([]logging.LoggingOption{})
	client, err := clientAuth(&authReq, emptyLogConfig)
	assert.NoError(t, err)
	assert.Equal(t, server.URL+"/api/v3", client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, "turbonomic-terraform-provider/1.1.0", client.Headers["User-Agent"])

}
