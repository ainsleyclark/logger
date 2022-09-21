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

package workplace

import (
	"bytes"
	"github.com/ainsleyclark/errors"
	mocks "github.com/ainsleyclark/logger/gen/mocks/test"
	"github.com/ainsleyclark/logger/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"os"
	"testing"
)

func TestNewHook(t *testing.T) {
	tt := map[string]struct {
		input string
		want  any
	}{
		"Success": {
			"token",
			nil,
		},
		"Error": {
			"",
			"token cannot be nil",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := NewHook(Options{Token: test.input})
			if err != nil {
				assert.Contains(t, err.Error(), test.want)
				return
			}
		})
	}
}

func TestHook_Fire(t *testing.T) {
	m := &mocks.Notifier{}
	m.On("Notify", mock.Anything, mock.Anything).
		Return(nil)
	h := Hook{wp: m}
	got := h.Fire(&logrus.Entry{})
	assert.Nil(t, got)
}

func TestHook_Process(t *testing.T) {
	entry := types.Entry{
		Message: "message",
		Data: map[string]any{
			"key":           "value",
			logrus.ErrorKey: "value",
		},
	}

	tt := map[string]struct {
		input     types.Entry
		mock      func(m *mocks.Notifier)
		formatter types.FormatMessageFunc
		want      any
	}{
		"Success": {
			entry,
			func(m *mocks.Notifier) {
				m.On("Notify", mock.Anything, mock.Anything).
					Return(nil)
			},
			nil,
			"",
		},
		"With Format": {
			entry,
			func(m *mocks.Notifier) {
				m.On("Notify", mock.Anything, mock.Anything).
					Return(nil)
			},
			func(entry types.Entry, args types.FormatMessageArgs) string {
				return "hello"
			},
			"",
		},
		"Error": {
			entry,
			func(m *mocks.Notifier) {
				m.On("Notify", mock.Anything, mock.Anything).
					Return(errors.New("error"))
			},
			nil,
			"error",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()
			m := &mocks.Notifier{}
			if test.mock != nil {
				test.mock(m)
			}
			h := Hook{
				wp:        m,
				options:   Options{FormatMessage: test.formatter},
				LogLevels: logrus.AllLevels,
			}
			h.process(test.input)
			assert.Contains(t, buf.String(), test.want)
		})
	}
}

func TestHook_FormatMessage(t *testing.T) {
	entry := types.Entry{
		Message: "message",
	}
	got := FormatMessage(entry, types.FormatMessageArgs{
		Service: "Service",
		Version: "v0.1.1",
		Prefix:  "PREFIX",
	})
	assert.Contains(t, got, "v0.1.1")
	assert.Contains(t, got, "Prefix")
}

func TestHook_Levels(t *testing.T) {
	h := Hook{LogLevels: logrus.AllLevels}
	want := logrus.AllLevels
	assert.Equal(t, want, h.Levels())
}
