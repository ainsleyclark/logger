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
)

func (t *LoggerTestSuite) TestDefaultHook_Fire() {
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
		"Dont Report": {
			&logrus.Entry{},
			defaultHook{
				config: &Config{
					report: func(e types.Entry) bool {
						return false
					},
				},
			},
			nil,
		},
		"OK": {
			&logrus.Entry{},
			defaultHook{
				wp: func(entry *logrus.Entry) error {
					return nil
				},
				mogrus: func(entry *logrus.Entry) error {
					return nil
				},
				config: &Config{
					report: types.DefaultReportFn,
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
					report: types.DefaultReportFn,
				},
			},
			"wp error",
		},
		"With Mogrus Error": {
			&logrus.Entry{},
			defaultHook{
				mogrus: func(entry *logrus.Entry) error {
					return errors.New("mogrus error")
				},
				config: &Config{
					report: types.DefaultReportFn,
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
