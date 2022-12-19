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

package logger

import (
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/logger/types"
	"github.com/sirupsen/logrus"
	"io"
)

func (t *LoggerTestSuite) TestDefaultHook_Fire() {
	L.SetOutput(io.Discard)

	tt := map[string]struct {
		input *logrus.Entry
		hook  defaultHook
		want  any
	}{
		"Nil": {
			nil,
			defaultHook{},
			nil,
		},
		"OK": {
			&logrus.Entry{},
			defaultHook{
				wp:     func(entry *logrus.Entry) error { return nil },
				mogrus: func(entry *logrus.Entry) error { return nil },
				slack:  func(entry *logrus.Entry) error { return nil },
				config: &Config{
					workplace: workplaceConfig{Report: types.DefaultReportFn},
					mongo:     mongoConfig{Report: types.DefaultReportFn},
					slack:     slackConfig{Report: types.DefaultReportFn},
				},
			},
			nil,
		},
		"Workplace Dont Report": {
			&logrus.Entry{},
			defaultHook{
				wp: func(entry *logrus.Entry) error {
					return nil
				},
				config: &Config{
					workplace: workplaceConfig{Report: func(e types.Entry) bool {
						return false
					}},
				},
			},
			nil,
		},
		"With Workplace Error": {
			&logrus.Entry{},
			defaultHook{
				wp: func(entry *logrus.Entry) error {
					return errors.New("wp error")
				},
				config: &Config{
					workplace: workplaceConfig{Report: types.DefaultReportFn},
				},
			},
			nil,
		},
		"Slack Dont Report": {
			&logrus.Entry{},
			defaultHook{
				slack: func(entry *logrus.Entry) error {
					return nil
				},
				config: &Config{
					slack: slackConfig{Report: func(e types.Entry) bool {
						return false
					}},
				},
			},
			nil,
		},
		"With Slack Error": {
			&logrus.Entry{},
			defaultHook{
				slack: func(entry *logrus.Entry) error {
					return errors.New("wp error")
				},
				config: &Config{
					slack: slackConfig{Report: types.DefaultReportFn},
				},
			},
			nil,
		},
		"Mogrus Dont Report": {
			&logrus.Entry{},
			defaultHook{
				mogrus: func(entry *logrus.Entry) error {
					return nil
				},
				config: &Config{
					mongo: mongoConfig{Report: func(e types.Entry) bool {
						return false
					}},
				},
			},
			nil,
		},
		"With Mogrus Error": {
			&logrus.Entry{},
			defaultHook{
				mogrus: func(entry *logrus.Entry) error {
					return errors.New("mogrus error")
				},
				config: &Config{
					mongo: mongoConfig{Report: types.DefaultReportFn},
				},
			},
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			err := test.hook.Fire(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.want, err)
		})
	}
}
