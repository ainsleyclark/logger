// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"github.com/ainsleyclark/errors"
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
	}
)

const (
	// DefaultPrefix is the default prefix used when none
	// is set.
	DefaultPrefix = "REDDICO"
	// DefaultStatus is the default status used when none
	// is set.
	DefaultStatus = "RED"
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
func (op *Options) WithMongoCollection(collection *mongo.Collection) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.mongoCollection = collection
	})
	return op
}

// WithWorkplaceNotifier sends errors that have been marked
// as errors.INTERNAL to a Workplace thread.
func (op *Options) WithWorkplaceNotifier(token, thread string) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.workplaceToken = token
		config.workplaceThread = thread
	})
	return op
}
