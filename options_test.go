// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
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
		"Database Name": {
			Config{
				service:     "service",
				mongoClient: &mongo.Client{},
			},
			"database name cannot be empty",
		},
		"Mongo Client": {
			Config{
				service:       "service",
				mongoDatabase: "database",
			},
			"mongo client cannot be nil",
		},
		"Workplace Thread": {
			Config{
				service:        "service",
				mongoClient:    &mongo.Client{},
				mongoDatabase:  "database",
				workplaceToken: "token",
			},
			"workplace thread cannot be nil",
		},
		"Workplace Token": {
			Config{
				service:         "service",
				mongoClient:     &mongo.Client{},
				mongoDatabase:   "database",
				workplaceThread: "thread",
			},
			"workplace token cannot be nil",
		},
		"Success": {
			Config{
				service:         "service",
				mongoClient:     &mongo.Client{},
				mongoDatabase:   "database",
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
		prefix:        DefaultPrefix,
		defaultStatus: DefaultStatus,
	}
	t.Equal(want, *got)
}

func (t *LoggerTestSuite) TestOptions() {
	opts := NewOptions().
		Service("service").
		Version("v0.0.1").
		DefaultStatus("status").
		Prefix("prefix").
		WithMongoClient(&mongo.Client{}, "database").
		WithWorkplaceNotifier("token", "thread")

	c := &Config{}
	for _, optFn := range opts.optFuncs {
		optFn(c)
	}

	t.Equal("service", c.service)
	t.Equal("v0.0.1", c.version)
	t.Equal("status", c.defaultStatus)
	t.Equal("prefix", c.prefix)
	t.Equal("database", c.mongoDatabase)
	t.Equal("token", c.workplaceToken)
	t.Equal("thread", c.workplaceThread)
}
