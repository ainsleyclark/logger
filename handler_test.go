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
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
			buf := t.Setup()

			handler := func(w http.ResponseWriter, r *http.Request) {
				test.input.Request = r
				w.WriteHeader(test.input.Status)
				Fire(test.input)
			}

			ts := httptest.NewServer(http.Handler(http.HandlerFunc(handler)))
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL+"/test", nil)
			assert.NoErrorf(t.T(), err, "error making new request")

			_, err = http.DefaultClient.Do(req)
			assert.NoErrorf(t.T(), err, "error performing request")

			t.Contains(buf.String(), test.want)
		})
	}
}
