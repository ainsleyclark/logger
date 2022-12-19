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
	"github.com/ainsleyclark/logger/types"
	"go.mongodb.org/mongo-driver/mongo"
)

func (t *LoggerTestSuite) TestConfig_Validate() {
	tt := map[string]struct {
		input Config
		want  any
	}{
		"Service": {
			Config{},
			"service name cannot be empty",
		},
		"Workplace Thread": {
			Config{
				service:   "service",
				workplace: workplaceConfig{Token: "token"},
			},
			"workplace thread cannot be nil",
		},
		"Workplace Token": {
			Config{
				service:   "service",
				workplace: workplaceConfig{Thread: "thread"},
			},
			"workplace token cannot be nil",
		},
		"Success": {
			Config{
				service:   "service",
				workplace: workplaceConfig{Token: "token", Thread: "thread"},
			},
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.Validate()
			if got != nil {
				t.Contains(got.Error(), test.want)
				return
			}
		})
	}
}

func (t *LoggerTestSuite) TestConfig_AssignDefaults() {
	c := Config{}
	got := c.assignDefaults()
	want := Config{
		prefix:        DefaultPrefix,
		defaultStatus: DefaultStatus,
		workplace:     workplaceConfig{Report: types.DefaultReportFn},
		mongo:         mongoConfig{Report: types.DefaultReportFn},
		slack:         slackConfig{Report: types.DefaultReportFn},
	}
	t.Equal(want.prefix, got.prefix)
	t.Equal(want.defaultStatus, got.defaultStatus)
	t.NotNil(got.workplace.Report)
	t.NotNil(got.mongo.Report)
	t.NotNil(got.slack.Report)
}

func (t *LoggerTestSuite) TestOptions() {
	opts := NewOptions().
		Service("service").
		Version("v0.0.1").
		DefaultStatus("status").
		Prefix("prefix").
		WithMongoCollection(&mongo.Collection{}, types.DefaultReportFn).
		WithWorkplaceNotifier("token", "thread", types.DefaultReportFn, nil).
		WithSlackNotifier("token", "channel", types.DefaultReportFn, nil)

	c := &Config{}
	for _, optFn := range opts.optFuncs {
		optFn(c)
	}

	t.Equal("service", c.service)
	t.Equal("v0.0.1", c.version)
	t.Equal("status", c.defaultStatus)
	t.Equal("prefix", c.prefix)
	t.Equal("token", c.workplace.Token)
	t.Equal("thread", c.workplace.Thread)
	t.Equal("token", c.slack.Token)
	t.Equal("channel", c.slack.Channel)
	t.NotNil(c.workplace.Report)
	t.NotNil(c.mongo.Report)
	t.NotNil(c.slack.Report)
}
