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

package types

import "github.com/sirupsen/logrus"

type (
	// Fields is an alias for logrus.Fields.
	Fields = logrus.Fields
	// Entry is an abstraction for logrus.Entry, which provides
	// helper functions.
	Entry logrus.Entry
	// ShouldReportFunc is the function used to determine if a
	// logger entry should be sent to a hook.
	ShouldReportFunc func(e Entry) bool
)

var (
	// DefaultReportFn is the default report function when
	// none is passed to the constructor.
	DefaultReportFn = func(e Entry) bool {
		return true
	}
)

// ToLogrusEntry transforms an Entry to logrus.Entry
func (e Entry) ToLogrusEntry() logrus.Entry {
	return logrus.Entry(e)
}

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

func (e *Entry) HasError() {

}
