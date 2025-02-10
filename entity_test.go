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

// The EntityTests struct is referenced from testdata.go which needs to be created based on testdata.go.template
package turboclient

import (
	"testing"
)

func convertEntityTest(entity EntityResults) TestEntity {
	return TestEntity{
		uuid:        entity.UUID,
		displayName: entity.DisplayName,
		className:   entity.ClassName}
}

func TestTurboEntity(t *testing.T) {

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
