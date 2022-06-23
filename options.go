// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
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
		mongoClient     *mongo.Client
		workplaceToken  string
		workplaceThread string //nolint
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

func (op *Options) Service(service string) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.service = service
	})
	return op
}

func (op *Options) WithMongoCollection(client *mongo.Client) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.mongoClient = client
	})
	return op
}

func (op *Options) WithWorkplaceNotifier(token string) *Options {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.workplaceToken = token
	})
	return op
}
