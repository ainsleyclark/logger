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
		L = logrus.New()
	}()
	L.ExitFunc = func(i int) {}
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
	t.Equal(buf, L.Out)
}

func (t *LoggerTestSuite) TestSetLevel() {
	defer func() {
		L = logrus.New()
	}()
	SetLevel(logrus.WarnLevel)
	t.Equal(logrus.WarnLevel, L.GetLevel())
}

func (t *LoggerTestSuite) TestSetLogger() {
	defer func() {
		L = logrus.New()
	}()
	l := L
	SetLogger(l)
	t.Equal(l, L)
}

//func (t *LoggerTestSuite) TestSetService() {
//
//	t.Run("Success", func() {
//		orig := config
//		defer func() {
//			config = orig
//		}()
//		SetService("service")
//		t.Equal("service", config.service)
//	})
//
//	//t.Run("Error", func() {
//	//	buf := t.Setup()
//	//	SetService("")
//	//	color.Greenln(b)
//	//})
//}
