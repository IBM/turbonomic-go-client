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
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEntity(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/GetEntity.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/entities/75941320319680", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	// Set the base URL for the client
	client.BaseURL = ts.URL

	// Create an EntityRequest with test parameters
	entityReq := EntityRequest{Uuid: "75941320319680"}

	// Call the GetEntity function
	entityResults, err := client.GetEntity(entityReq)

	// Assert that the function returned the expected result and no error
	assert.NoError(t, err)
	assert.Equal(t, "75941320319680", entityResults.UUID)
	assert.Equal(t, "test-vm", entityResults.DisplayName)
	assert.Equal(t, "VirtualMachine", entityResults.ClassName)
}

func TestTurboEntityIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests, to run set environment variable INTEGRATION")
	}
	newClientOpts := ClientParameters{Hostname: TurboHost, Username: TurboUser, Password: TurboPass, Skipverify: DoNotVerify}

	c, err := NewClient(&newClientOpts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
		t.FailNow()
	}

	for _, tt := range EntityTests {
		entity, err := c.GetEntity(EntityRequest{Uuid: tt.uuid})

		if err != nil {
			t.Errorf("error: %s", err.Error())
			t.FailNow()
		}
		expectedResult := convertEntityTest(*entity)
		if tt != expectedResult {
			t.Errorf("error: got %s expected %s", tt, expectedResult)
		}

	}
}

func convertEntityTest(entity EntityResults) TestEntity {
	return TestEntity{
		uuid:        entity.UUID,
		displayName: entity.DisplayName,
		className:   entity.ClassName}
}
