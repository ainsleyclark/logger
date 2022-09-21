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

import "net/http"

func (t *LoggerTestSuite) TestEntry_IsHTTP() {
	tt := map[string]struct {
		input Fields
		want  bool
	}{
		"True": {
			Fields{
				"status_code": http.StatusOK,
				"client_ip":   "127.0.0.1",
				"request_url": "https://github.com/ainsleyclark/logger",
			},
			true,
		},
		"False": {
			Fields{},
			false,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			e := Entry{Data: test.input}
			t.Equal(test.want, e.IsHTTP())
		})
	}
}
