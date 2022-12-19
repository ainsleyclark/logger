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

package slack

import (
	"bytes"
	"errors"
	"github.com/ainsleyclark/logger/types"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestNewHook(t *testing.T) {
	got := NewHook(Options{})
	assert.NotNil(t, got.sendFunc)
	assert.NotNil(t, got.options)
	assert.Equal(t, logrus.AllLevels, got.LogLevels)
}

func TestHook_Fire(t *testing.T) {
	h := Hook{sendFunc: func(channelID string, options ...slack.MsgOption) (string, string, error) {
		return "", "", nil
	}}
	got := h.Fire(&logrus.Entry{})
	assert.Nil(t, got)
}

func TestHook_Process(t *testing.T) {
	entry := types.Entry{
		Message: "message",
		Data: map[string]any{
			"key":          "value",
			types.ErrorKey: "value",
		},
	}

	tt := map[string]struct {
		mock sendSlackFunc
		want any
	}{
		"Success": {
			func(channelID string, options ...slack.MsgOption) (string, string, error) {
				return "", "", nil
			},
			"",
		},
		"Error": {
			func(channelID string, options ...slack.MsgOption) (string, string, error) {
				return "", "", errors.New("error")
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
			h := Hook{
				sendFunc:  test.mock,
				LogLevels: logrus.AllLevels,
			}
			h.process(entry)
			assert.Contains(t, buf.String(), test.want)
		})
	}
}

func TestHook_Levels(t *testing.T) {
	h := Hook{LogLevels: logrus.AllLevels}
	want := logrus.AllLevels
	assert.Equal(t, want, h.Levels())
}
