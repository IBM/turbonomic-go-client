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

// The ActionTests struct is referenced from testdata.go which needs to be created based on testdata.go.template
package turboclient

import (
	"testing"
)

func TestGetActions(t *testing.T) {

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
			t.Errorf("received %d entities, wanted %d entites", len(respAction), 1)
			t.FailNow()
		}
		if respAction[0].Target.DisplayName != tt.displayName {
			t.Errorf("error: got %s expected %s", respAction[0].Target.DisplayName, tt.displayName)
		}
	}

}
