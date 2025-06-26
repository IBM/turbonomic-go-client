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
	"fmt"
)

// Parameters for searching Turbonomic's API
type SearchRequest struct {
	Name             string
	EntityType       string
	EnvironmentType  string
	CloudType        string
	CaseSensitive    bool
	CommonReqParams  CommonReqParams
	SearchParameters map[string]string
}

// Criterion for a Turbonomic API search request
type Criteria struct {
	CaseSensitive bool   `json:"caseSensitive"`
	ExpType       string `json:"expType"`
	ExpVal        string `json:"expVal"`
	FilterType    string `json:"filterType"`
}

// Body for POST request of Turbonomic API search request
type SearchDTO struct {
	CriteriaList    []Criteria `json:"criteriaList"`
	LogicalOperator string     `json:"logicalOperator"`
	ClassName       string     `json:"className"`
	Scope           string     `json:"scope,omitempty"`
	EnvironmentType string     `json:"environmentType,omitempty"`
	CloudType       string     `json:"cloudType,omitempty"`
}

// Results of a search request to Turbonomic's API
type SearchResults []struct {
	UUID            string `json:"uuid"`
	DisplayName     string `json:"displayName"`
	ClassName       string `json:"className"`
	EnvironmentType string `json:"environmentType"`
	DiscoveredBy    struct {
		UUID        string `json:"uuid"`
		DisplayName string `json:"displayName"`
		Category    string `json:"category"`
		Type        string `json:"type"`
		Readonly    bool   `json:"readonly"`
	} `json:"discoveredBy"`
	VendorIds         struct{} `json:"vendorIds"`
	State             string   `json:"state"`
	Severity          string   `json:"severity"`
	CostPrice         float64  `json:"costPrice"`
	SeverityBreakdown struct{} `json:"severityBreakdown"`
	Template          struct {
		Price       float64 `json:"price"`
		Discovered  bool    `json:"discovered"`
		EnableMatch bool    `json:"enableMatch"`
		DisplayName string  `json:"displayName"`
	} `json:"template"`
	Aspects struct {
		VirtualMachineAspect struct {
			Os                string   `json:"os"`
			IP                []string `json:"ip"`
			NumVCPUs          int      `json:"numVCPUs"`
			EbsOptimized      bool     `json:"ebsOptimized"`
			ResourceID        string   `json:"resourceId"`
			CreationTimeStamp int      `json:"creationTimeStamp"`
			Type              string   `json:"type"`
		} `json:"virtualMachineAspect"`
		VirtualDisksAspect struct {
			VirtualDisks []struct {
				UUID        string `json:"uuid"`
				DisplayName string `json:"displayName"`
				Tier        string `json:"tier"`
				Stats       []struct {
					Name     string `json:"name"`
					Capacity struct {
						Max   int `json:"max"`
						Min   int `json:"min"`
						Avg   int `json:"avg"`
						Total int `json:"total"`
					} `json:"capacity"`
					Filters []struct {
						Type        string      `json:"type"`
						Value       string      `json:"value"`
						DisplayName interface{} `json:"displayName"`
					} `json:"filters"`
					Units  string `json:"units"`
					Values struct {
						Max   int `json:"max"`
						Min   int `json:"min"`
						Avg   int `json:"avg"`
						Total int `json:"total"`
					} `json:"values"`
					Value int `json:"value"`
				} `json:"stats"`
				AttachedVirtualMachine struct {
					UUID        string `json:"uuid"`
					DisplayName string `json:"displayName"`
					ClassName   string `json:"className"`
				} `json:"attachedVirtualMachine"`
				Provider struct {
					UUID        string `json:"uuid"`
					DisplayName string `json:"displayName"`
					ClassName   string `json:"className"`
				} `json:"provider"`
				DataCenter struct {
					UUID        string `json:"uuid"`
					DisplayName string `json:"displayName"`
					ClassName   string `json:"className"`
				} `json:"dataCenter"`
				EnvironmentType string `json:"environmentType"`
				LastModified    int64  `json:"lastModified"`
				BusinessAccount struct {
					UUID            string `json:"uuid"`
					DisplayName     string `json:"displayName"`
					ClassName       string `json:"className"`
					EnvironmentType string `json:"environmentType"`
					DiscoveredBy    struct {
						UUID        string `json:"uuid"`
						DisplayName string `json:"displayName"`
						Category    string `json:"category"`
						Type        string `json:"type"`
						Readonly    bool   `json:"readonly"`
					} `json:"discoveredBy"`
					VendorIds struct {
						Turbonomicamp string `json:"turbonomicamp"`
					} `json:"vendorIds"`
					State             string `json:"state"`
					Severity          string `json:"severity"`
					SeverityBreakdown struct {
						NORMAL int `json:"NORMAL"`
					} `json:"severityBreakdown"`
					Tags struct {
						Usage []string `json:"Usage"`
					} `json:"tags"`
					Staleness string `json:"staleness"`
				} `json:"businessAccount"`
				SnapshotID        string  `json:"snapshotId"`
				Encryption        string  `json:"encryption"`
				AttachmentState   string  `json:"attachmentState"`
				HourlyBilledOps   float64 `json:"hourlyBilledOps"`
				CreationTimeStamp int64   `json:"creationTimeStamp"`
				ResourceID        string  `json:"resourceId"`
			} `json:"virtualDisks"`
			Type string `json:"type"`
		} `json:"virtualDisksAspect"`
	} `json:"aspects"`
	Tags struct{} `json:"tags"`
}

// Retrives the results of a search of Turbonomic's API based on provided parameters
func (c *Client) SearchEntityByName(searchReq SearchRequest) (SearchResults, error) {

	// var filterType string
	filterType, err := c.getFilterType(searchReq.EntityType)
	if err != nil {
		return SearchResults{}, err
	}

	searchCriteria := SearchDTO{
		EnvironmentType: searchReq.EnvironmentType,
		CloudType:       searchReq.CloudType,
		LogicalOperator: "OR",
		ClassName:       searchReq.EntityType,
		CriteriaList: []Criteria{
			{
				CaseSensitive: searchReq.CaseSensitive,
				ExpType:       "EQ",
				ExpVal:        searchReq.Name,
				FilterType:    filterType,
			},
		},
	}

	return c.SearchEntities(searchCriteria, searchReq.CommonReqParams)
}

func (c *Client) SearchEntities(
	searchCriteria SearchDTO, reqParams CommonReqParams) (SearchResults, error) {

	dtoBuf := new(bytes.Buffer)
	if err := json.NewEncoder(dtoBuf).Encode(searchCriteria); err != nil {
		return nil, err
	}

	restResp, err := c.request(RequestOptions{Method: "POST", Path: "/search", ReqDTO: dtoBuf,
		CommonReqParams: CommonReqParams{
			Headers:         reqParams.Headers,
			QueryParameters: reqParams.QueryParameters}})
	if err != nil {
		return nil, err
	}

	var searchResults SearchResults
	if err := json.Unmarshal(restResp, &searchResults); err != nil {
		return nil, err
	}

	return searchResults, err
}

// Helper function to enable the use of entity type as the filter instead of
// longer parameter names required by Turbonomic's API
func (c *Client) getFilterType(entityType string) (string, error) {
	entityMap := map[string]string{
		"VirtualMachine": "vmsByName",
		"VirtualVolume":  "virtualVolumeByName",
		"DatabaseServer": "databaseByName",
	}
	filterType := entityMap[entityType]
	if filterType != "" {
		return filterType, nil
	}
	return "", fmt.Errorf("entity type of %s not supported", entityType)
}
