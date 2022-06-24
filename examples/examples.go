// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package examples

import (
	"github.com/ainsleyclark/errors"
	"github.com/krang-backlink/logger"
	"net/http"
	"time"
)

// BasicUsage godoc
func BasicUsage() {
	logger.Trace("Trace Entry")
	logger.Debug("Debug Entry")
	logger.Info("Info Entry")
	logger.Warn("Warn Entry")
	logger.Error("Error Entry")
	logger.Fatal("Fatal Entry")
	logger.Panic("Panic Entry")
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
		Request:      r,
		Status:       http.StatusOK,
		Message:      "Message from API",
		Data:         map[string]any{},
		RequestTime:  time.Now(),
		ResponseTime: time.Now(),
		Latency:      100,
	})
}
