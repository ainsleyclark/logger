// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package notify

import (
	"bytes"
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/mogrus"
	"github.com/krang-backlink/logger/gen/mocks/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"os"
	"testing"
)

func TestNewFireHook(t *testing.T) {
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
			"Error creating Workplace Client",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			m := &mocks.Notifier{}
			m.On("Notify", mock.Anything, mock.Anything).
				Return(nil)
			hook, err := NewFireHook(test.input)
			if err != nil {
				assert.Contains(t, errors.Message(err), test.want)
				return
			}
			hook(mogrus.Entry{})
		})
	}
}

func TestFire(t *testing.T) {
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
			fire(m, test.input)
			assert.Contains(t, buf.String(), test.want)
		})
	}
}
