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

// This is a sample file that needs to be populated in order to run integration tests.
// All tests are pulling data from this file.

package turboclient

var TurboHost string = ""
var TurboUser string = ""
var TurboPass string = ""
var DoNotVerify bool = true

type TestEntity struct {
	uuid        string
	displayName string
	className   string
}

type TestSearch struct {
	entityName      string
	entityType      string
	environmentType string
	caseSensitive   bool
	queryParameters map[string]string
	uuid            string
}

type TestAction struct {
	uuid         string
	actionStates []string
	actionTypes  []string
	displayName  string
}

type TestStats struct {
	uuid        string
	endDate     string
	statistics  []StatisticRequest
	displayName string
}

var EntityTests = []TestEntity{
	{},
	{},
	{},
}

var SearchTests = []TestSearch{
	{"", "VirtualMachine", "ONPREM", true, map[string]string{"query_type": "EXACT"}, ""}, // Added to avoid linter errors in build
	{},
}

var ActionTests = []TestAction{
	{},
	{},
}

var StatsTests = []TestStats{
	{},
}
