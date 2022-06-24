// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/logger/internal/stdout"
	"github.com/ainsleyclark/logger/internal/workplace"
	"github.com/ainsleyclark/mogrus"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func (t *LoggerTestSuite) TestNew() {
	tt := map[string]struct {
		input  func() *Options
		mogrus func(ctx context.Context, opts mogrus.Options) (logrus.Hook, error)
		hook   func(opts workplace.Options) (*workplace.Hook, error)
		want   any
	}{
		"Simple": {
			func() *Options {
				return NewOptions().Service("service")
			},
			mogrus.New,
			workplace.NewHook,
			nil,
		},
		"Validation Error": {
			func() *Options {
				return NewOptions()
			},
			mogrus.New,
			workplace.NewHook,
			"service name cannot be empty",
		},
		"With Workplace": {
			func() *Options {
				return NewOptions().Service("service").WithWorkplaceNotifier("token", "thread")
			},
			mogrus.New,
			workplace.NewHook,
			nil,
		},
		"Workplace Error": {
			func() *Options {
				return NewOptions().Service("service").WithWorkplaceNotifier("token", "thread")
			},
			mogrus.New,
			func(opts workplace.Options) (*workplace.Hook, error) {
				return nil, errors.New("hook error")
			},
			"hook error",
		},
		"With Mogrus": {
			func() *Options {
				return NewOptions().Service("service").WithMongoCollection(&mongo.Collection{})
			},
			func(ctx context.Context, opts mogrus.Options) (logrus.Hook, error) {
				return &stdout.Hook{}, nil
			},
			workplace.NewHook,
			nil,
		},
		"Mogrus Error": {
			func() *Options {
				return NewOptions().Service("service").WithMongoCollection(&mongo.Collection{})
			},
			func(ctx context.Context, opts mogrus.Options) (logrus.Hook, error) {
				return nil, errors.New("mogrus error")
			},
			workplace.NewHook,
			"mogrus error",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			origMogrus := newMogrus
			origHook := newHook
			defer func() {
				newMogrus = origMogrus
				newHook = origHook
			}()
			newMogrus = test.mogrus
			newHook = test.hook
			err := New(context.TODO(), test.input())
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
		})
	}
}

func (t *LoggerTestSuite) TestLogger() {
	tt := map[string]struct {
		fn   func()
		want string
	}{
		"Trace": {
			func() {
				Trace("trace")
			},
			"trace",
		},
		"Debug": {
			func() {
				Debug("debug")
			},
			"debug",
		},
		"Info": {
			func() {
				Info("info")
			},
			"info",
		},
		"Warn": {
			func() {
				Warn("warning")
			},
			"warning",
		},
		"Error": {
			func() {
				Error("error")
			},
			"error",
		},
		"With Field": {
			func() {
				WithField("test", "with-field").Error()
			},
			"with-field",
		},
		"With Fields": {
			func() {
				WithFields(logrus.Fields{"test": "with-fields"}).Error()
			},
			"with-fields",
		},
		"With Error": {
			func() {
				WithError(&errors.Error{Code: "code", Message: "message", Operation: "op", Err: fmt.Errorf("err")}).Error()
			},
			"[code] code [msg] message [op] op [error] err",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			buf := t.Setup()
			test.fn()
			t.Contains(buf.String(), test.want)
		})
	}
}

func (t *LoggerTestSuite) TestLogger_Fatal() {
	buf := t.Setup() // nolint
	defer func() {
		logger = logrus.New()
	}()
	logger.ExitFunc = func(i int) {}
	Fatal("fatal")
	t.Contains(buf.String(), "fatal")
}

func (t *LoggerTestSuite) TestLogger_Panic() {
	buf := t.Setup()
	t.Panics(func() {
		Panic("panic")
	})
	t.Contains(buf.String(), "panic")
}

func (t *LoggerTestSuite) TestLogger_SetOutput() {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	t.Equal(buf, logger.Out)
}

func (t *LoggerTestSuite) TestSetLevel() {
	defer func() {
		logger = logrus.New()
	}()
	SetLevel(logrus.WarnLevel)
	t.Equal(logrus.WarnLevel, logger.GetLevel())
}

func (t *LoggerTestSuite) TestSetLogger() {
	defer func() {
		logger = logrus.New()
	}()
	l := logger
	SetLogger(l)
	t.Equal(l, logger)
}
