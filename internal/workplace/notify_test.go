// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workplace

import (
	"bytes"
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/mogrus"
	mocks "github.com/krang-backlink/logger/gen/mocks/test"
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
	h := Hook{}
	got := h.Fire(&logrus.Entry{})
	assert.Nil(t, got)
}

func TestHook_Process(t *testing.T) {
	entry := mogrus.Entry{
		Level:   "info",
		Message: "message",
		Data: map[string]any{
			"key":           "value",
			logrus.ErrorKey: "value",
		},
		Error: &mogrus.Error{Code: errors.INTERNAL},
	}

	tt := map[string]struct {
		input mogrus.Entry
		mock  func(m *mocks.Notifier)
		want  any
	}{
		"Success": {
			entry,
			func(m *mocks.Notifier) {
				m.On("Notify", mock.Anything, mock.Anything).
					Return(nil)
			},
			"",
		},
		"Nil Error": {
			mogrus.Entry{Level: "info", Message: "message"},
			nil,
			"",
		},
		"Not Internal": {
			mogrus.Entry{Level: "info", Message: "message", Error: &mogrus.Error{Code: errors.INVALID}},
			nil,
			"",
		},
		"Error": {
			entry,
			func(m *mocks.Notifier) {
				m.On("Notify", mock.Anything, mock.Anything).
					Return(errors.New("error"))
			},
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
				options:   Options{},
				LogLevels: logrus.AllLevels,
			}

			h.process(test.input)

			assert.Contains(t, buf.String(), test.want)
		})
	}
}

func TestHook_Levels(t *testing.T) {
	h := Hook{LogLevels: logrus.AllLevels}
	want := logrus.AllLevels
	assert.Equal(t, want, h.Levels())
}
