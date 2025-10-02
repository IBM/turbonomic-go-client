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
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
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
		t.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/search", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "{\"criteriaList\":[],\"logicalOperator\":\"AND\","+
			"\"className\":\"VirtualMachine\",\"scope\":\"null\",\"environmentType\":\"ONPREM\"}\n",
			string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write header: %s", err.Error())
			t.FailNow()
		}
	}))
	defer server.Close()

	// Setup search criteria and request parameters
	searchCriteria := SearchDTO{CriteriaList: []Criteria{}, LogicalOperator: "AND", ClassName: "VirtualMachine", Scope: "null", EnvironmentType: "ONPREM"}
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
		t.Fatal("Error when opening file: ", err)
	}

	tests := []struct {
		searchReq    SearchRequest
		responseBody string
	}{
		{
			SearchRequest{
				Name:            "test-vm",
				EntityType:      "VirtualMachine",
				EnvironmentType: "ONPREM",
				CaseSensitive:   true},
			"{\"criteriaList\":[" +
				"{" +
				"\"caseSensitive\":true," +
				"\"expType\":\"EQ\"," +
				"\"expVal\":\"test-vm\"," +
				"\"filterType\":\"vmsByName\"" +
				"}]," +
				"\"logicalOperator\":\"AND\"," +
				"\"className\":\"VirtualMachine\"," +
				"\"environmentType\":\"ONPREM\"" +
				"}\n",
		},
		{
			SearchRequest{
				Name:            "test-vm",
				EntityType:      "VirtualMachine",
				EnvironmentType: "CLOUD",
				CaseSensitive:   true,
				OSNames:         []string{"Linux", "Windows"}},
			"{\"criteriaList\":[" +
				"{" +
				"\"caseSensitive\":true," +
				"\"expType\":\"EQ\"," +
				"\"expVal\":\"test-vm\"," +
				"\"filterType\":\"vmsByName\"" +
				"},{" +
				"\"caseSensitive\":false," +
				"\"expType\":\"RXEQ\"," +
				"\"expVal\":\"Linux.*|Windows.*\"," +
				"\"filterType\":\"vmsByGuestName\"" +
				"}]," +
				"\"logicalOperator\":\"AND\"," +
				"\"className\":\"VirtualMachine\"," +
				"\"environmentType\":\"CLOUD\"" +
				"}\n",
		},
	}

	for _, test := range tests {

		// Create a test server with the mock response
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/search", r.URL.Path)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, test.responseBody, string(body))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(mockResponse)); err != nil {
				t.Errorf("failed to create client: %s", err.Error())
				t.FailNow()
			}
		}))
		defer ts.Close()

		// Set the base URL for the client
		client.BaseURL = ts.URL

		// Call the GetEntity function
		searchResults, err := client.SearchEntityByName(test.searchReq)

		// Assert that the function returned the expected result and no error
		assert.NoError(t, err)
		assert.Equal(t, "75928674621956", searchResults[0].UUID)
		assert.Equal(t, "test-vm", searchResults[0].DisplayName)
		assert.Equal(t, "VirtualMachine", searchResults[0].ClassName)
		assert.Equal(t, "ONPREM", searchResults[0].EnvironmentType)
		assert.Equal(t, "vm-121", searchResults[0].VendorIds["vc1.turbo.com"])
		assert.Equal(t, 1, len(searchResults))
	}
}

func TestSearchEntityByVendorId(t *testing.T) {
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
		t.Fatal("Error when opening file: ", err)
	}

	tests := []struct {
		searchReq    SearchRequestByVendorId
		responseBody string
	}{
		{
			SearchRequestByVendorId{
				EntityType: "VirtualMachine",
				VendorId:   "vm-121"},
			"{\"criteriaList\":[" +
				"{" +
				"\"caseSensitive\":false," +
				"\"expType\":\"RXEQ\"," +
				"\"expVal\":\"vm-121\"," +
				"\"filterType\":\"vmsByVendorId\"" +
				"}]," +
				"\"logicalOperator\":\"AND\"," +
				"\"className\":\"VirtualMachine\"" +
				"}\n",
		},
	}

	for _, test := range tests {
		// Create a test server with the mock response
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/search", r.URL.Path)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, test.responseBody, string(body))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(mockResponse)); err != nil {
				t.Errorf("failed to create client: %s", err.Error())
				t.FailNow()
			}
		}))
		defer ts.Close()

		// Set the base URL for the client
		client.BaseURL = ts.URL

		// Call the SearchEntityByName function
		searchResults, err := client.SearchEntityByVendorId(test.searchReq)

		// Assert that the function returned the expected result and no error
		assert.NoError(t, err)
		assert.Equal(t, "75928674621956", searchResults[0].UUID)
		assert.Equal(t, "vm-121", searchResults[0].VendorIds["vc1.turbo.com"])
		assert.Equal(t, "VirtualMachine", searchResults[0].ClassName)
		assert.Equal(t, "ONPREM", searchResults[0].EnvironmentType)
		assert.Equal(t, 1, len(searchResults))
	}
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
			CloudType:       tt.cloudType,
			OSNames:         tt.osName,
		}

		ans, err := c.SearchEntityByName(searchReq)
		if err != nil {
			t.Errorf("error: %s", err.Error())
			t.FailNow()
		}
		if len(ans) != len(tt.uuids) {
			t.Errorf("received %d entities, wanted %d entites for entity name %s", len(ans), len(tt.uuids), tt.entityName)
			t.FailNow()
		}
		for _, result := range ans {
			if !slices.Contains(tt.uuids, result.UUID) {
				t.Errorf("received %s, wanted %s", result.UUID, tt.uuids)
			}
		}
	}
}

func TestSearchEntityByVendorIdIntegration(t *testing.T) {
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
		// Skip tests that don't have a vendor ID
		if tt.vendorId == "" {
			continue
		}

		searchReq := SearchRequestByVendorId{
			EntityType: tt.entityType,
			VendorId:   tt.vendorId,
		}

		ans, err := c.SearchEntityByVendorId(searchReq)
		if err != nil {
			t.Errorf("error: %s", err.Error())
			t.FailNow()
		}
		if len(ans) != len(tt.uuids) {
			t.Errorf("received %d entities, wanted %d entities for vendor ID %s", len(ans), len(tt.uuids), tt.vendorId)
			t.FailNow()
		}
		for _, result := range ans {
			if !slices.Contains(tt.uuids, result.UUID) {
				t.Errorf("received %s, wanted %s", result.UUID, tt.uuids)
			}
		}
	}
}
func TestSearchVolumeById(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/SearchVolumeById.json")
	if err != nil {
		t.Fatal("Error when opening file: ", err)
	}

	tests := []struct {
		searchReq    SearchRequestByVendorId
		responseBody string
	}{
		{
			SearchRequestByVendorId{
				EntityType: "VirtualVolume",
				VendorId:   "vol-030ba747139817215"},
			"{\"criteriaList\":[" +
				"{" +
				"\"caseSensitive\":false," +
				"\"expType\":\"RXEQ\"," +
				"\"expVal\":\"vol-030ba747139817215\"," +
				"\"filterType\":\"volumeById\"" +
				"}]," +
				"\"logicalOperator\":\"AND\"," +
				"\"className\":\"VirtualVolume\"" +
				"}\n",
		},
	}

	for _, test := range tests {
		// Create a test server with the mock response
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/search", r.URL.Path)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, test.responseBody, string(body))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(mockResponse)); err != nil {
				t.Errorf("failed to create client: %s", err.Error())
				t.FailNow()
			}
		}))
		defer ts.Close()

		// Set the base URL for the client
		client.BaseURL = ts.URL

		// Call the SearchEntityByVendorId function
		searchResults, err := client.SearchEntityByVendorId(test.searchReq)

		// Assert that the function returned the expected result and no error
		assert.NoError(t, err)
		assert.Equal(t, "76153836730880", searchResults[0].UUID)
		assert.Equal(t, "vol-030ba747139817215", searchResults[0].DisplayName)
		assert.Equal(t, "VirtualVolume", searchResults[0].ClassName)
		assert.Equal(t, "CLOUD", searchResults[0].EnvironmentType)
		assert.Equal(t, "vol-030ba747139817215", searchResults[0].VendorIds["aws-test"])
		assert.Equal(t, 1, len(searchResults))
	}
}
