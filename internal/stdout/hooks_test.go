// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdout

import (
	"bytes"
	"fmt"
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

type mockFormatErr struct{}

func (m *mockFormatErr) Format(entry *logrus.Entry) ([]byte, error) {
	return nil, fmt.Errorf("err")
}

type mockFormat struct{}

func (m *mockFormat) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte("test"), nil
}

type mockWriterErr struct{}

func (m *mockWriterErr) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("err")
}

// SetupHooks is a helper function for setting up
// the hooks for testing.
func SetupHooks(writer io.Writer) Hook {
	return Hook{
		Writer: writer,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	}
}

func TestHook_Fire(t *testing.T) {
	buf := &bytes.Buffer{}

	tt := map[string]struct {
		input io.Writer
		entry *logrus.Entry
		want  any
	}{
		"Error Entry": {
			&bytes.Buffer{},
			&logrus.Entry{
				Logger: &logrus.Logger{Formatter: &mockFormatErr{}},
			},
			"Error obtaining the entry string",
		},
		"Error Writer": {
			&mockWriterErr{},
			&logrus.Entry{
				Logger: &logrus.Logger{Formatter: &mockFormat{}},
			},
			"Error writing entry to io.Writer",
		},
		"Success": {
			buf,
			&logrus.Entry{
				Logger: &logrus.Logger{Formatter: &mockFormat{}},
			},
			"test",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			h := SetupHooks(test.input)
			err := h.Fire(test.entry)
			if err != nil {
				assert.Contains(t, errors.Message(err), test.want)
				return
			}
			assert.Equal(t, test.want, buf.String())
		})
	}
}

func TestHook_Levels(t *testing.T) {
	h := SetupHooks(nil)
	want := []logrus.Level{
		logrus.InfoLevel,
		logrus.DebugLevel,
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
	assert.Equal(t, want, h.Levels())
}
