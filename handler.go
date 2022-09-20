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
		L.WithFields(fields).Info(f.Message)
		return
	}

	L.WithFields(fields).Error(f.Message)
}
