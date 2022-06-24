// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"fmt"
	"github.com/ainsleyclark/errors"
	"github.com/sirupsen/logrus"
	"time"
)

func (t *LoggerTestSuite) TestFormatter() {
	var (
		now       = time.Now()
		nowStr    = now.Format(time.StampMilli)
		prefix    = "[TEST]"
		defStatus = "TEST"
	)

	tt := map[string]struct {
		entry *logrus.Entry
		want  string
	}{
		"Debug": {
			&logrus.Entry{
				Level:   logrus.DebugLevel,
				Message: "message",
			},
			fmt.Sprintf(prefix+" %s | %s | [DEBUG] | [msg] message\n", nowStr, defStatus),
		},
		"Info": {
			&logrus.Entry{
				Level:   logrus.InfoLevel,
				Message: "message",
			},
			fmt.Sprintf(prefix+" %s | %s | [INFO]  | [msg] message\n", nowStr, defStatus),
		},
		"Warning": {
			&logrus.Entry{
				Level:   logrus.WarnLevel,
				Message: "message",
			},
			fmt.Sprintf(prefix+" %s | %s | [WARNING] | [msg] message\n", nowStr, defStatus),
		},
		"Error": {
			&logrus.Entry{
				Level:   logrus.ErrorLevel,
				Message: "message",
			},
			fmt.Sprintf(prefix+" %s | %s | [ERROR] | [msg] message\n", nowStr, defStatus),
		},
		"Fatal": {
			&logrus.Entry{
				Level:   logrus.FatalLevel,
				Message: "message",
			},
			fmt.Sprintf(prefix+" %s | %s | [FATAL] | [msg] message\n", nowStr, defStatus),
		},
		"Panic": {
			&logrus.Entry{
				Level:   logrus.PanicLevel,
				Message: "message",
			},
			fmt.Sprintf(prefix+" %s | %s | [PANIC] | [msg] message\n", nowStr, defStatus),
		},
		"Fields": {
			&logrus.Entry{
				Data: logrus.Fields{
					"fields": logrus.Fields{"key1": "test1"},
				},
				Level: logrus.InfoLevel,
			},
			fmt.Sprintf(prefix+" %s | %s | [INFO]  | key1: test1\n", nowStr, defStatus),
		},
		"Print Error Pointer": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": &errors.Error{Code: "INTERNAL", Message: "message", Operation: "operation", Err: fmt.Errorf("error")},
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(prefix+" %s | %s | [ERROR] | [code] INTERNAL [msg] message [op] operation [error] error\n", nowStr, defStatus),
		},
		"Print Error Non Pointer": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": errors.Error{Code: "INTERNAL", Message: "message", Operation: "operation", Err: fmt.Errorf("error")},
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(prefix+" %s | %s | [ERROR] | [code] INTERNAL [msg] message [op] operation [error] error\n", nowStr, defStatus),
		},
		"Nil To Error": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": 1,
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(prefix+" %s | %s | [ERROR]\n", nowStr, defStatus),
		},
		"Print Error": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": fmt.Errorf("error"),
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(prefix+" %s | %s | [ERROR] | [error] error\n", nowStr, defStatus),
		},
		"Print Error String": {
			&logrus.Entry{
				Data: logrus.Fields{
					"error": "error",
				},
				Level: logrus.ErrorLevel,
			},
			fmt.Sprintf(prefix+" %s | %s | [ERROR] | [error] error\n", nowStr, defStatus),
		},
		"Server Success": {
			&logrus.Entry{
				Data: logrus.Fields{
					"status_code":    200,
					"client_ip":      "127.0.0.1",
					"request_method": "GET",
					"request_url":    "/page",
					"data_length":    0,
				},
				Level: logrus.InfoLevel,
			},
			fmt.Sprintf(prefix+" %s | 200 | [INFO]  | 127.0.0.1 |   GET    \"/page\"\n", nowStr),
		},
		"Server Not Found": {
			&logrus.Entry{
				Data: logrus.Fields{
					"status_code":    404,
					"client_ip":      "127.0.0.1",
					"request_method": "GET",
					"request_url":    "/page",
					"data_length":    0,
				},
				Level: logrus.InfoLevel,
			},
			fmt.Sprintf(prefix+" %s | 404 | [INFO]  | 127.0.0.1 |   GET    \"/page\"\n", nowStr),
		},
		"Message": {
			&logrus.Entry{
				Data: logrus.Fields{
					"message": "message",
				},
				Level: logrus.InfoLevel,
			},
			fmt.Sprintf(prefix+" %s | %s | [INFO]  | [msg] message\n", nowStr, defStatus),
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			test.entry.Time = now
			f := formatter{
				Config: &Config{
					prefix:        "test",
					defaultStatus: "test",
				},
				Colours: false,
			}
			got, err := f.Format(test.entry)
			t.NoError(err)
			t.Equal(test.want, string(got))
		})
	}
}
