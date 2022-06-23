// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"github.com/ainsleyclark/errors"
	"github.com/go-chi/chi/v5"
	"github.com/krang-backlink/api/delivery/api/shared/httptest"
	"net/http"
)

func (t *LoggerTestSuite) TestHandler() {
	tt := map[string]struct {
		input FireHook
		want  any
	}{
		"200": {
			FireHook{
				Status: http.StatusOK,
			},
			"200 | [INFO]",
		},
		"500": {
			FireHook{
				Status: http.StatusInternalServerError,
			},
			"500 | [ERROR]",
		},
		"Message": {
			FireHook{
				Status:  http.StatusOK,
				Message: "message",
			},
			"[msg] message",
		},
		"Error": {
			FireHook{
				Status: http.StatusOK,
				Data:   errors.NewInternal(errors.New("error"), "message", "op"),
			},
			"[code] internal [msg] message [op] op [error] error",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			s := httptest.New(t.T())
			buf := t.Setup()
			r := chi.NewRouter()
			r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
				test.input.Request = r
				w.WriteHeader(test.input.Status)
				Fire(test.input)
			})
			s.RequestAndServe(http.MethodGet, "/test", nil, r)
			t.Contains(buf.String(), test.want)
		})
	}
}
