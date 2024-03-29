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

func TestHook_FormatMessage(t *testing.T) {
	entry := Entry{
		Message: "message",
		Data: Fields{
			ErrorKey: errors.NewInternal(errors.New("error"), "message", "op"),
			FieldKey: map[string]any{
				"test": "hello",
			},
		},
	}
	got := DefaultFormatMessageFn(entry, FormatMessageArgs{
		Service: "Service",
		Version: "v0.1.1",
		Prefix:  "PREFIX",
	})
	assert.Contains(t, got, "v0.1.1")
	assert.Contains(t, got, "Prefix")
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

func TestEntry_Fields(t *testing.T) {
	var empty logrus.Fields

	tt := map[string]struct {
		input Entry
		want  any
	}{
		"Nil": {
			Entry{},
			empty,
		},
		"Bad Cast": {
			Entry{
				Data: map[string]any{
					FieldKey: errors.New("error"),
				},
			},
			empty,
		},
		"OK": {
			Entry{
				Data: map[string]any{
					FieldKey: map[string]any{"test": "hello"},
				},
			},
			Fields{"test": "hello"},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.Fields()
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
					ErrorKey: nil,
				},
			},
			false,
		},
		"With Error": {
			Entry{
				Data: map[string]any{
					ErrorKey: errors.New("error"),
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

func TestEntry_Error(t *testing.T) {
	var empty *errors.Error

	tt := map[string]struct {
		input Entry
		want  any
	}{
		"Nil": {
			Entry{},
			empty,
		},
		"With Error": {
			Entry{
				Data: map[string]any{
					ErrorKey: errors.New("error"),
				},
			},
			"error",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			got := test.input.Error()
			if got != nil {
				assert.Contains(t, got.Error(), test.want)
				return
			}
			assert.Equal(t, test.want, got)
		})
	}
}
