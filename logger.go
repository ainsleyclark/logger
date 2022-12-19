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
	"context"
	"github.com/ainsleyclark/logger/internal/hooks/stdout"
	"github.com/ainsleyclark/logger/types"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var (
	// L is an alias for the standard logrus Logger.
	L = logrus.New()
)

// New creates a new standard L and sets logging levels
// dependent on environment variables.
func New(ctx context.Context, opts ...*Options) error {
	c := &Config{}
	for _, opt := range opts {
		for _, optFn := range opt.optFuncs {
			optFn(c)
		}
	}
	err := c.Validate()
	if err != nil {
		return err
	}
	return initialise(ctx, c.assignDefaults())
}

// Trace logs a trace message with args.
func Trace(args ...any) {
	L.Trace(args...)
}

// Debug logs a debug message with args.
func Debug(args ...any) {
	L.Debug(args...)
}

// Info logs ab info message with args.
func Info(args ...any) {
	L.Info(args...)
}

// Warn logs a warn message with args.
func Warn(args ...any) {
	L.Warn(args...)
}

// Error logs an error message with args.
func Error(args ...any) {
	L.Error(args...)
}

// Fatal logs a fatal message with args.
func Fatal(args ...any) {
	L.Fatal(args...)
}

// Panic logs a panic message with args.
func Panic(args ...any) {
	L.Panic(args...)
}

// WithField logs with field, sets a new map containing
// "fields".
func WithField(key string, value any) *logrus.Entry {
	return L.WithFields(logrus.Fields{types.FieldKey: logrus.Fields{
		key: value,
	}})
}

// WithFields logs with fields, sets a new map containing
// "fields".
func WithFields(fields types.Fields) *logrus.Entry {
	return L.WithFields(logrus.Fields{types.FieldKey: fields})
}

// WithError - Logs with a custom error.
func WithError(err any) *logrus.Entry {
	return L.WithField(types.ErrorKey, err)
}

// SetOutput sets the output of the L to an io.Writer,
// useful for testing.
func SetOutput(writer io.Writer) {
	L.SetOutput(writer)
}

// SetLevel sets the level of the L.
func SetLevel(level logrus.Level) {
	L.SetLevel(level)
}

// SetLogger sets the application L.
func SetLogger(l *logrus.Logger) {
	L = l
}

// initialise sets the standard log level, sets the
// log formatter and discards the stdout.
func initialise(ctx context.Context, cfg *Config) error { //nolint
	L.SetLevel(logrus.TraceLevel)

	L.SetFormatter(&formatter{
		Config:          cfg,
		TimestampFormat: "2006-01-02 15:04:05",
		Colours:         true,
	})

	// Send all logs to nowhere by default.
	L.SetOutput(io.Discard)

	// Send logs with level higher than warning to stderr.
	L.AddHook(&stdout.Hook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	// Send info and debug logs to stdout.
	L.AddHook(&stdout.Hook{
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.TraceLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})

	// Add the WP & Mogrus hooks to the logger.
	err := addHooks(ctx, cfg)
	if err != nil {
		return err
	}

	return nil
}
