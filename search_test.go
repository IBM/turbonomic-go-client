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
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchEntities(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/SearchAllVms.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/search", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "{\"criteriaList\":[],\"logicalOperator\":\"AND\","+
			"\"className\":\"VirtualMachine\",\"scope\":\"null\",\"environmentType\":\"\"}\n",
			string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Setup search criteria and request parameters
	searchCriteria := SearchDTO{CriteriaList: []Criteria{}, LogicalOperator: "AND", ClassName: "VirtualMachine", Scope: "null"}
	reqParams := CommonReqParams{}

	// Set the base URL for the client
	client.BaseURL = server.URL

	// Call the function being tested
	searchResults, err := client.SearchEntities(searchCriteria, reqParams)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 20, len(searchResults))
	assert.Equal(t, "75941320319539", searchResults[4].UUID)
	assert.Equal(t, "test-vm-06", searchResults[5].DisplayName)
	assert.Equal(t, "VirtualMachine", searchResults[7].ClassName)
	assert.Equal(t, "CLOUD", searchResults[15].EnvironmentType)
}

func TestSearchEntityByName(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/SearchEntityByName.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/search", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "{\"criteriaList\":[{\"caseSensitive\":true,\"expType\":\"EQ\","+
			"\"expVal\":\"test-vm\",\"filterType\":\"vmsByName\"}],"+
			"\"logicalOperator\":\"OR\",\"className\":\"VirtualMachine\","+
			"\"environmentType\":\"ONPREM\"}\n",
			string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	// Set the base URL for the client
	client.BaseURL = ts.URL

	// Create an EntityRequest with test parameters
	searchReq := SearchRequest{
		Name:            "test-vm",
		EntityType:      "VirtualMachine",
		EnvironmentType: "ONPREM",
		CaseSensitive:   true,
	}

	// Call the GetEntity function
	searchResults, err := client.SearchEntityByName(searchReq)

	// Assert that the function returned the expected result and no error
	assert.NoError(t, err)
	assert.Equal(t, "75928674621956", searchResults[0].UUID)
	assert.Equal(t, "test-vm", searchResults[0].DisplayName)
	assert.Equal(t, "VirtualMachine", searchResults[0].ClassName)
	assert.Equal(t, "ONPREM", searchResults[0].EnvironmentType)
	assert.Equal(t, 1, len(searchResults))
}

func TestSearchEntitiesIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests, to run set environment variable INTEGRATION")
	}
	newClientOpts := ClientParameters{Hostname: TurboHost, Username: TurboUser, Password: TurboPass, Skipverify: DoNotVerify}
	c, err := NewClient(&newClientOpts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
		t.FailNow()
	}

	for _, tt := range SearchTests {
		searchReq := SearchRequest{
			Name:            tt.entityName,
			EntityType:      tt.entityType,
			EnvironmentType: tt.environmentType,
			CaseSensitive:   tt.caseSensitive,
		}

		ans, err := c.SearchEntityByName(searchReq)
		if err != nil {
			t.Errorf("error: %s", err.Error())
			t.FailNow()
		}
		if len(ans) != 1 {
			t.Errorf("received %d entities, wanted %d entites for entity name %s", len(ans), 1, tt.entityName)
			t.FailNow()
		}
		if ans[0].UUID != tt.uuid {
			t.Errorf("received %s, wanted %s", ans[0].UUID, tt.uuid)
		}

	}

}
