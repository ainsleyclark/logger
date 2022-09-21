// Copyright 2022 Ainsley Clark. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDefaultReportFn(t *testing.T) {
	got := DefaultReportFn(Entry{})
	assert.Equal(t, true, got)
}

func TestEntry_ToLogrusEntry(t *testing.T) {
	e := Entry{}
	got := e.ToLogrusEntry()
	assert.IsType(t, logrus.Entry{}, got)
}

func TestEntry_IsHTTP(t *testing.T) {
	tt := map[string]struct {
		input Entry
		want  bool
	}{
		"True": {
			Entry{
				Data: map[string]any{
					"status_code": http.StatusOK,
					"client_ip":   "127.0.0.1",
					"request_url": "https://github.com/ainsleyclark/logger",
				},
			},
			true,
		},
		"False": {
			Entry{},
			false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.IsHTTP()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestEntry_HasError(t *testing.T) {
	tt := map[string]struct {
		input Entry
		want  bool
	}{
		"Nil": {
			Entry{},
			false,
		},
		"No Error": {
			Entry{
				Data: map[string]any{
					logrus.ErrorKey: nil,
				},
			},
			false,
		},
		"With Error": {
			Entry{
				Data: map[string]any{
					logrus.ErrorKey: errors.New("error"),
				},
			},
			true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.HasError()
			assert.Equal(t, test.want, got)
		})
	}
}
