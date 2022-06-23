// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// FireHook represents the data needed to log out a message
// or data for http Requests.
type FireHook struct {
	Request      *http.Request
	Status       int
	Message      string
	Data         any
	RequestTime  time.Time
	ResponseTime time.Time
	Latency      float64
}

// Fire fires a FireHook to Logrus from a http request.
func Fire(f FireHook) {
	endTime := time.Now()
	latency := time.Since(f.RequestTime)

	var err *errors.Error
	e, ok := f.Data.(*errors.Error)
	if ok {
		err = e
	}

	fields := logrus.Fields{
		"status_code":    f.Status,
		"latency_time":   endTime.Sub(f.RequestTime),
		"client_ip":      f.Request.RemoteAddr,
		"request_method": f.Request.Method,
		"request_url":    f.Request.RequestURI,
		"referer":        f.Request.Referer(),
		"user_agent":     f.Request.UserAgent(),
		"request_time":   f.RequestTime,
		"response_time":  f.ResponseTime,
		"duration":       float64(latency.Nanoseconds()) / float64(1000),
		"message":        f.Message,
		logrus.ErrorKey:  err,
	}

	if f.Status >= 200 && f.Status < 300 {
		logger.WithFields(fields).Info()
		return
	}

	logger.WithFields(fields).Error()
}
