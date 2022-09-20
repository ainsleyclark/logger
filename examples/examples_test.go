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

package examples

import (
	"context"
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/logger"
	"net/http"
	"time"
)

// QuickStart godoc
func QuickStart() error {
	err := logger.New(context.TODO(), logger.NewOptions().Service("service"))
	if err != nil {
		return err
	}

	logger.Trace("Trace Entry")
	logger.Debug("Debug Entry")
	logger.Info("Info Entry")
	logger.Warn("Warn Entry")
	logger.Error("Error Entry")
	logger.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error()
	logger.Fatal("Fatal Entry")

	return nil
}

// Fields allow you to log out key value pairs to the logger
// that will appear under data. The simplest way to use the
// logger is simply the package-level exported logger.
func Fields() {
	logger.WithFields(logger.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")
}

// WithError godoc
func WithError() {
	logger.WithError(errors.NewInternal(errors.New("error"), "message", "op")).Error()
}

// Middleware is provided out of the box in the form of a fire hook.
// Upon receiving a request from the API, calling `logger.Fire`
// will send the log entry to stdout with detailed request
// information and meta.
func Middleware(r *http.Request) {
	logger.Fire(logger.FireHook{
		Request: r,
		Status:  http.StatusOK,
		Message: "Message from API",
		Data: map[string]any{
			"key": "value",
		},
		RequestTime:  time.Now(),
		ResponseTime: time.Now(),
		Latency:      100,
	})
}
