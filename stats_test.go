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

// For Integrations tests, the StatsTests struct is referenced from testdata.go which needs to be created based on testdata.go.template.
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

func TestClient_GetStats(t *testing.T) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := &Client{
		BaseURL: "/api/v3",
		HTTPClient: &http.Client{
			Transport: customTransport,
		},
	}

	// Mock response from the Turbonomic API
	mockResponse, err := os.ReadFile("./testfiles/GetStats.json")
	if err != nil {
		t.Fatal("Error when opening file: ", err)
	}

	// Create a test server with the mock response
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/stats/76097737480945", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		// Check for expected request body content
		assert.Contains(t, string(body), "\"endDate\":\"+10m\"")
		assert.Contains(t, string(body), "\"name\":\"StorageAccess\"")
		assert.Contains(t, string(body), "\"relatedEntityType\":\"VirtualVolume\"")

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

	// Create a StatsRequest with test parameters
	statsReq := StatsRequest{
		EntityUUID: "76097737480945",
		EndDate:    "+10m",
		Statistics: []StatisticRequest{
			{
				Name:              "StorageAccess",
				RelatedEntityType: "VirtualVolume",
				Filters: []Filter{
					{
						Type:  "relation",
						Value: "sold",
					},
				},
			},
			{
				Name:              "StorageAmount",
				RelatedEntityType: "VirtualVolume",
				Filters: []Filter{
					{
						Type:  "relation",
						Value: "sold",
					},
				},
			},
			{
				Name:              "IOThroughput",
				RelatedEntityType: "VirtualVolume",
				Filters: []Filter{
					{
						Type:  "relation",
						Value: "sold",
					},
				},
			},
		},
	}

	// Call the GetStats function
	statsResp, err := client.GetStats(statsReq)

	// Assert that the function returned the expected result and no error
	assert.NoError(t, err)
	assert.NotEmpty(t, statsResp)
	assert.Equal(t, "vol-05f7c906f860b4d3c", statsResp[0].DisplayName)
	assert.NotEmpty(t, statsResp[0].Statistics)
	assert.Equal(t, "StorageAccess", statsResp[0].Statistics[0].Name)
}

func TestGetStatsIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests, to run set environment variable INTEGRATION")
	}

	newClientOpts := ClientParameters{Hostname: TurboHost, Username: TurboUser, Password: TurboPass, Skipverify: DoNotVerify}
	c, err := NewClient(&newClientOpts)
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
		t.FailNow()
	}

	for _, tt := range StatsTests {
		statsReq := StatsRequest{
			EntityUUID: tt.uuid,
			EndDate:    tt.endDate,
			Statistics: tt.statistics,
		}

		statsResp, err := c.GetStats(statsReq)
		if err != nil {
			t.Errorf("error: %s", err.Error())
			t.FailNow()
		}

		if len(statsResp) == 0 {
			t.Errorf("received empty response for entity %s", tt.uuid)
			t.FailNow()
		}

		if statsResp[0].DisplayName != tt.displayName {
			t.Errorf("error: got %s expected %s", statsResp[0].DisplayName, tt.displayName)
		}
	}
}
