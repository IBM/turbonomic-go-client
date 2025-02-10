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
	"bytes"
	"encoding/json"
	"log/slog"
)

// Parameters for retriving an entity from Turbonomic's API
type EntityRequest struct {
	Uuid             string
	CommonReqOptions CommonReqParams
}

// Results from GetEntity request
type EntityResults struct {
	UUID            string `json:"uuid,omitempty"`
	DisplayName     string `json:"displayName,omitempty"`
	ClassName       string `json:"className,omitempty"`
	EnvironmentType string `json:"environmentType,omitempty"`
	DiscoveredBy    struct {
		UUID        string `json:"uuid,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		Category    string `json:"category,omitempty"`
		Type        string `json:"type,omitempty"`
		Readonly    bool   `json:"readonly,omitempty"`
	} `json:"discoveredBy,omitempty"`
	VendorIds struct {
		Turbonomicamp string `json:"turbonomicamp,omitempty"`
	} `json:"vendorIds,omitempty"`
	State             string  `json:"state,omitempty"`
	Severity          string  `json:"severity,omitempty"`
	CostPrice         float64 `json:"costPrice,omitempty"`
	SeverityBreakdown struct {
		UNKNOWN  int `json:"UNKNOWN,omitempty"`
		NORMAL   int `json:"NORMAL,omitempty"`
		MINOR    int `json:"MINOR,omitempty"`
		MAJOR    int `json:"MAJOR,omitempty"`
		CRITICAL int `json:"CRITICAL,omitempty"`
	} `json:"severityBreakdown,omitempty"`
	Providers []struct {
		UUID        string `json:"uuid,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		ClassName   string `json:"className,omitempty"`
	} `json:"providers,omitempty"`
	Template struct {
		UUID        string  `json:"uuid,omitempty"`
		DisplayName string  `json:"displayName,omitempty"`
		Price       float64 `json:"price,omitempty"`
		Discovered  bool    `json:"discovered,omitempty"`
		EnableMatch bool    `json:"enableMatch,omitempty"`
	} `json:"template,omitempty"`
	Tags      struct{} `json:"tags,omitempty"`
	Staleness string   `json:"staleness,omitempty"`
}

// Retrives entity based on its provided uuid
func (c *Client) GetEntity(reqOpts EntityRequest) (*EntityResults, error) {

	restResp, err := c.request(RequestOptions{Method: "GET", Path: "/entities/" + reqOpts.Uuid, ReqDTO: new(bytes.Buffer),
		CommonReqParams: CommonReqParams{
			Headers:         reqOpts.CommonReqOptions.Headers,
			QueryParameters: reqOpts.CommonReqOptions.QueryParameters}})

	if err != nil {
		return nil, err
	}
	slog.Debug(string(restResp))

	var entityResults EntityResults
	json.Unmarshal(restResp, &entityResults)

	return &entityResults, err
}
