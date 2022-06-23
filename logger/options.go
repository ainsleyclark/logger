// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// Options defines the configuration needed for creating
	// a new Logger.
	Options struct {
		// Prefix is the string written to the log before any
		// message.
		Prefix string
		// DefaultStatus is the status code for HTTP requests
		// when none is set.
		DefaultStatus string
	}
	Config struct {
		prefix          string
		defaultStatus   string
		service         string
		mongoClient     *mongo.Client
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

// assignDefaults the default prefix and status in the case
// when there is none set.
func (o *Options) assignDefaults() Options {
	if o.Prefix == "" {
		o.Prefix = DefaultPrefix
	}
	if o.DefaultStatus == "" {
		o.DefaultStatus = DefaultStatus
	}
	return *o
}

// optionFunc is a function type that configures a T instance.
type optionFunc func(config *Config)

// Options is the type used to configure a new T instance.
type TestOptions struct {
	optFuncs []optionFunc
}

// NewOptions creates an empty Options instance.
func NewOptions() *TestOptions {
	return &TestOptions{}
}

func (op *TestOptions) Prefix(prefix string) *TestOptions {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.prefix = prefix
	})
	return op
}

func (op *TestOptions) DefaultStatus(status string) *TestOptions {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.defaultStatus = status
	})
	return op
}

func (op *TestOptions) Service(service string) *TestOptions {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.service = service
	})
	return op
}

func (op *TestOptions) WithMongoClient(client *mongo.Client) *TestOptions {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.mongoClient = client
	})
	return op
}

func (op *TestOptions) WithWorkplaceNotifier(token string) *TestOptions {
	op.optFuncs = append(op.optFuncs, func(config *Config) {
		config.workplaceToken = token
	})
	return op
}
