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
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// Config defines the configuration needed for creating
	// a new Logger.
	Config struct {
		version         string
		prefix          string
		defaultStatus   string
		service         string
		mongoCollection *mongo.Collection
		workplaceToken  string
		workplaceThread string
		workplaceReport types.ShouldReportFunc
		mongoReport     types.ShouldReportFunc
	}
)

const (
	// DefaultPrefix is the default prefix used when none
	// is set.
	DefaultPrefix = "LOGGER"
	// DefaultStatus is the default status used when none
	// is set.
	DefaultStatus = "LOG"
)

// Validate ensures the configuration is sanity checked
// before creating a new Logger.
func (c *Config) Validate() error {
	if c.service == "" {
		return errors.New("service name cannot be empty")
	}
	if c.workplaceToken != "" && c.workplaceThread == "" {
		return errors.New("workplace thread cannot be nil")
	}
	if c.workplaceToken == "" && c.workplaceThread != "" {
		return errors.New("workplace token cannot be nil")
	}
	return nil
}

// assignDefaults the default prefix and status in the case
// when there is none set.
func (c *Config) assignDefaults() *Config {
	if c.prefix == "" {
		c.prefix = DefaultPrefix
	}
	if c.defaultStatus == "" {
		c.defaultStatus = DefaultStatus
	}
	if c.workplaceReport == nil {
		c.workplaceReport = types.DefaultReportFn
	}
	if c.mongoReport == nil {
		c.mongoReport = types.DefaultReportFn
	}
	return c
}

// optionFunc is a function type that configures a config instance.
type optionFunc func(config *Config)

// Options is the type used to configure a new config instance.
type Options struct {
	optFuncs []optionFunc
}

// NewOptions creates an empty Options instance.
func NewOptions() *Options {
	return &Options{}
}

// Version is the currently running version of the service
// or application.
func (op *Options) Version(version string) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.version = version
	})
	return op
}

// Prefix is the string written to the log before any
// message.
func (op *Options) Prefix(prefix string) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.prefix = prefix
	})
	return op
}

// DefaultStatus is the status code for HTTP requests
// when none is set.
func (op *Options) DefaultStatus(status string) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.defaultStatus = status
	})
	return op
}

// Service is used for Mongo logging, and general stdout logs.
// This name will correlate to the name of the Mongo
// database, if WithMongoCollection is called.
func (op *Options) Service(service string) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.service = service
	})
	return op
}

// WithMongoCollection allows for logging directly to Mongo.
func (op *Options) WithMongoCollection(collection *mongo.Collection, fn types.ShouldReportFunc) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.mongoCollection = collection
		config.mongoReport = fn
	})
	return op
}

// WithWorkplaceNotifier sends errors that have been marked
// as errors.INTERNAL to a Workplace thread.
func (op *Options) WithWorkplaceNotifier(token, thread string, fn types.ShouldReportFunc) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.workplaceToken = token
		config.workplaceThread = thread
		config.workplaceReport = fn
	})
	return op
}
