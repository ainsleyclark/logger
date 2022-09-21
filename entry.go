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

import "github.com/sirupsen/logrus"

type (
	// Entry is an abstraction for logrus.Entry, which provides
	// helper functions.
	Entry logrus.Entry
)

// IsHTTP returns true if the fields contain HTTP
// key and value pairs.
func (e *Entry) IsHTTP() bool {
	_, status := e.Data["status_code"]
	_, ip := e.Data["client_ip"]
	_, url := e.Data["request_url"]
	if status && ip && url {
		return true
	}
	return false
}
