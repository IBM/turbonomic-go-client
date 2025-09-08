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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetActionsByUUID(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/GetActionsByUuid.json")
	if err != nil {
		t.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/entities/75941320319680/actions", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "{\"actionStateList\":[\"READY\"],\"actionTypeList\":[\"RESIZE\"],\"detailLevel\":\"EXECUTION\"}\n", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Fail()
			t.Log(err)
		}
	}))
	defer ts.Close()

	// Set the base URL for the client
	client.BaseURL = ts.URL

	// Create an ActionsRequest with test parameters
	actionReq := ActionsRequest{
		Uuid:        "75941320319680",
		ActionState: []string{"READY"},
		ActionType:  []string{"RESIZE"},
		DetailLevel: "EXECUTION",
	}

	// Call the GetActionsByUUID function
	actionResults, err := client.GetActionsByUUID(actionReq)

	// Assert that the function returned the expected result and no error
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actionResults))
	assert.Equal(t, "638911097668880", actionResults[0].UUID)
	assert.Equal(t, "75941320319680", actionResults[0].Target.UUID)
}

func TestGetActionsByUUIDCompound(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/GetActionsByUuidCompound.json")
	if err != nil {
		t.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/entities/76084922964120/actions", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "{\"actionStateList\":[\"READY\"],\"actionTypeList\":[\"RESIZE\"],\"detailLevel\":\"EXECUTION\"}\n", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Fail()
			t.Log(err)
		}
	}))
	defer ts.Close()

	// Set the base URL for the client
	client.BaseURL = ts.URL

	// Create an ActionsRequest with test parameters
	actionReq := ActionsRequest{
		Uuid:        "76084922964120",
		ActionState: []string{"READY"},
		ActionType:  []string{"RESIZE"},
		DetailLevel: "EXECUTION",
	}

	// Call the GetActionsByUUID function
	actionResults, err := client.GetActionsByUUID(actionReq)

	// Assert that the function returned the expected result and no error
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actionResults))
	assert.Equal(t, 3, len(actionResults[0].CompoundActions))
	assert.Equal(t, "639054921252942", actionResults[0].UUID)
	assert.Equal(t, "76084922964120", actionResults[0].Target.UUID)
	assert.Equal(t, "75878878480048", actionResults[0].CompoundActions[0].Target.DiscoveredBy.UUID)
	assert.Equal(t, "262144.0", actionResults[0].CompoundActions[0].CurrentValue)
	assert.Equal(t, "491520.0", actionResults[0].CompoundActions[0].NewValue)
	assert.Equal(t, "1048576.0", actionResults[0].CompoundActions[1].CurrentValue)
	assert.Equal(t, "1310720.0", actionResults[0].CompoundActions[1].NewValue)
	assert.Equal(t, "200.0", actionResults[0].CompoundActions[2].CurrentValue)
	assert.Equal(t, "10.0", actionResults[0].CompoundActions[2].NewValue)

}

func TestGetActionsByUUIDMulti(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/GetActionsByUuidMulti.json")
	if err != nil {
		t.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/entities/75930461864800/actions", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		assert.Equal(t, "{\"actionStateList\":[\"READY\"],\"actionTypeList\":[\"RESIZE\"],\"detailLevel\":\"EXECUTION\"}\n", string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Fail()
			t.Log(err)
		}
	}))
	defer ts.Close()

	// Set the base URL for the client
	client.BaseURL = ts.URL

	// Create an ActionsRequest with test parameters
	actionReq := ActionsRequest{
		Uuid:        "75930461864800",
		ActionState: []string{"READY"},
		ActionType:  []string{"RESIZE"},
		DetailLevel: "EXECUTION",
	}

	// Call the GetActionsByUUID function
	actionResults, err := client.GetActionsByUUID(actionReq)

	// Assert that the function returned the expected result and no error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(actionResults))
	assert.Equal(t, "638883006725506", actionResults[0].UUID)
	assert.Equal(t, "75930461864800", actionResults[0].Target.UUID)
	assert.Equal(t, "2621440.0", actionResults[0].CurrentValue)
	assert.Equal(t, "3670016.0", actionResults[0].NewValue)
	assert.Equal(t, "638929431495649", actionResults[1].UUID)
	assert.Equal(t, "1.0", actionResults[1].CurrentValue)
	assert.Equal(t, "2.0", actionResults[1].NewValue)

}

func TestGetActionsIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests, to run set environment variable INTEGRATION")
	}
	newClientOpts := ClientParameters{Hostname: TurboHost, Username: TurboUser, Password: TurboPass, Skipverify: DoNotVerify}
	c, err := NewClient(&newClientOpts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
		t.FailNow()
	}

	for _, tt := range ActionTests {
		respAction, err := c.GetActionsByUUID(ActionsRequest{Uuid: tt.uuid, ActionState: tt.actionStates, ActionType: tt.actionTypes})

		if err != nil {
			t.Errorf("error: %s", err.Error())
			t.FailNow()
		}
		if len(respAction) != 1 {
			t.Errorf("received %d actions, wanted %d actions -> %s", len(respAction), 1, tt.displayName)
			t.FailNow()
		}
		if respAction[0].Target.DisplayName != tt.displayName {
			t.Errorf("error: got %s expected %s", respAction[0].Target.DisplayName, tt.displayName)
		}
	}
}
