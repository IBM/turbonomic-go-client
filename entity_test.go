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
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/IBM/turbonomic-go-client/logging"
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
		Logger: logging.NewSlogLogger(),
		Ctx:    context.Background(),
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/GetEntity.json")
	if err != nil {
		t.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/entities/75941320319680", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write header: %s", err.Error())
			t.FailNow()
		}
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

func TestTagEntity(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
		Logger: logging.NewSlogLogger(),
		Ctx:    context.Background(),
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/TagEntity.json")
	if err != nil {
		t.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/entities/75941320319680/tags", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write header: %s", err.Error())
			t.FailNow()
		}
	}))
	defer ts.Close()

	// Set the base URL for the client
	client.BaseURL = ts.URL

	// Create an TagEntityRequest with test parameters
	tagEntityReq := TagEntityRequest{Uuid: "75941320319680", Tags: []Tag{{Key: "Turbo_Team", Values: []string{"AppInfra_Integrations"}}}}

	// Call the TagEntity function
	tagEntityResults, err := client.TagEntity(tagEntityReq)

	// Assert that the function returned the expected result and no error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(tagEntityResults))
	assert.Equal(t, "Turbo_Team", tagEntityResults[0].Key)
	assert.Equal(t, "AppInfra_Integrations", tagEntityResults[0].Values[0])
	assert.Equal(t, "AppInfra_Integrations A", tagEntityResults[0].Values[1])
	assert.Equal(t, "Turbo_Owner", tagEntityResults[1].Key)
	assert.Equal(t, "Turbonomic_Appinfra_Integrations", tagEntityResults[1].Values[0])
}

func TestGetEntityTags(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
		Logger: logging.NewSlogLogger(),
		Ctx:    context.Background(),
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/TagEntity.json")
	if err != nil {
		t.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/entities/75941320319680/tags", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write header: %s", err.Error())
			t.FailNow()
		}
	}))
	defer ts.Close()

	// Set the base URL for the client
	client.BaseURL = ts.URL

	// Create an EntityRequest with test parameters
	entityReq := EntityRequest{Uuid: "75941320319680"}

	// Call the GetEntityTags function
	entityTagsResults, err := client.GetEntityTags(entityReq)

	// Assert that the function returned the expected result and no error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entityTagsResults))
	assert.Equal(t, "Turbo_Team", entityTagsResults[0].Key)
	assert.Equal(t, "AppInfra_Integrations", entityTagsResults[0].Values[0])
	assert.Equal(t, "AppInfra_Integrations A", entityTagsResults[0].Values[1])
	assert.Equal(t, "Turbo_Owner", entityTagsResults[1].Key)
	assert.Equal(t, "Turbonomic_Appinfra_Integrations", entityTagsResults[1].Values[0])
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
