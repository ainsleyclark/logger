// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
)

var (
	// logger is an alias for the standard logger.
	logger = logrus.New()
)

type (
	// Fields is an alias for logrus.Fields.
	Fields = logrus.Fields
)

// New creates a new standard logger and sets logging levels
// dependent on environment variables.
func New(ctx context.Context, opts ...*Options) error {
	c := &Config{}
	for _, opt := range opts {
		for _, optFn := range opt.optFuncs {
			optFn(c)
		}
	}
	return initialise(ctx, c.assignDefaults())
}

// Trace logs a trace message with args.
func Trace(args ...any) {
	logger.Trace(args...)
}

// Debug logs a debug message with args.
func Debug(args ...any) {
	logger.Debug(args...)
}

// Info logs ab info message with args.
func Info(args ...any) {
	logger.Info(args...)
}

// Warn logs a warn message with args.
func Warn(args ...any) {
	logger.Warn(args...)
}

// Error logs an error message with args.
func Error(args ...any) {
	logger.Error(args...)
}

// Fatal logs a fatal message with args.
func Fatal(args ...any) {
	logger.Fatal(args...)
}

// Panic logs a panic message with args.
func Panic(args ...any) {
	logger.Panic(args...)
}

// WithField logs with field, sets a new map containing
// "fields".
func WithField(key string, value any) *logrus.Entry {
	return logger.WithFields(logrus.Fields{"fields": logrus.Fields{
		key: value,
	}})
}

// WithFields logs with fields, sets a new map containing
// "fields".
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(logrus.Fields{"fields": fields})
}

// WithError - Logs with a custom error.
func WithError(err any) *logrus.Entry {
	return logger.WithField(logrus.ErrorKey, err)
}

// SetOutput sets the output of the logger to an io.Writer,
// useful for testing.
func SetOutput(writer io.Writer) {
	logger.SetOutput(writer)
}

// SetLevel sets the level of the logger.
func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

// SetLogger sets the application logger.
func SetLogger(l *logrus.Logger) {
	logger = l
}

// initialise sets the standard log level, sets the
// log formatter and discards the stdout.
func initialise(ctx context.Context, cfg *Config) error {
	logger.SetLevel(logrus.TraceLevel)

	logger.SetFormatter(&Formatter{
		Config:          cfg,
		TimestampFormat: "2006-01-02 15:04:05",
		Colours:         true,
	})

	// Send all logs to nowhere by default.
	logger.SetOutput(ioutil.Discard)

	// Send logs with level higher than warning to stderr.
	logger.AddHook(&WriterHook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	// Send info and debug logs to stdout
	logger.AddHook(&WriterHook{
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.TraceLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
		},
	})

	if cfg.mongoClient == nil {
		return nil
	}

	//hook, err := mogrus.New(ctx, mogrus.Options{
	//	Collection: collection,
	//	FireHook:   hook,
	//	ExpirationLevels: mogrus.ExpirationLevels{
	//		logrus.TraceLevel: time.Hour * 24,
	//		logrus.DebugLevel: time.Hour * 24,
	//		logrus.InfoLevel:  time.Hour * 24 * 7,
	//		logrus.ErrorLevel: time.Hour * 24 * 7 * 4,
	//		logrus.WarnLevel:  time.Hour * 24 * 7 * 4,
	//		logrus.PanicLevel: time.Hour * 24 * 7 * 4 * 6,
	//		logrus.FatalLevel: time.Hour * 24 * 7 * 4 * 6,
	//	},
	//})
	//if err != nil {
	//	return err
	//}
	//logger.AddHook(hook)
	//
	//// Add the Mogrus hook if a collection is passed.
	//if collection != nil {
	//
	//}

	return nil
}
