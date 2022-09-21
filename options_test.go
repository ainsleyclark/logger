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
				service:        "service",
				workplaceToken: "token",
			},
			"workplace thread cannot be nil",
		},
		"Workplace Token": {
			Config{
				service:         "service",
				workplaceThread: "thread",
			},
			"workplace token cannot be nil",
		},
		"Success": {
			Config{
				service:         "service",
				workplaceToken:  "token",
				workplaceThread: "thread",
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
		prefix:          DefaultPrefix,
		defaultStatus:   DefaultStatus,
		workplaceReport: types.DefaultReportFn,
		mongoReport:     types.DefaultReportFn,
	}
	t.Equal(want.prefix, got.prefix)
	t.Equal(want.defaultStatus, got.defaultStatus)
	t.NotNil(got.workplaceReport)
	t.NotNil(got.mongoReport)
}

func (t *LoggerTestSuite) TestOptions() {
	opts := NewOptions().
		Service("service").
		Version("v0.0.1").
		DefaultStatus("status").
		Prefix("prefix").
		WithMongoCollection(&mongo.Collection{}, types.DefaultReportFn).
		WithWorkplaceNotifier("token", "thread", types.DefaultReportFn, nil)

	c := &Config{}
	for _, optFn := range opts.optFuncs {
		optFn(c)
	}

	t.Equal("service", c.service)
	t.Equal("v0.0.1", c.version)
	t.Equal("status", c.defaultStatus)
	t.Equal("prefix", c.prefix)
	t.Equal("token", c.workplaceToken)
	t.Equal("thread", c.workplaceThread)
	t.NotNil(c.workplaceReport)
	t.NotNil(c.mongoReport)
}
