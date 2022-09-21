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

import (
	"bytes"
	"fmt"
	"github.com/ainsleyclark/errors"
	"github.com/enescakir/emoji"
	"github.com/sirupsen/logrus"
	"strings"
)

type (
	// Fields is an alias for logrus.Fields.
	Fields = logrus.Fields
	// Entry is an abstraction for logrus.Entry, which provides
	// helper functions.
	Entry logrus.Entry
	// ShouldReportFunc is the function used to determine if a
	// logger entry should be sent to a hook.
	ShouldReportFunc func(e Entry) bool
	// FormatMessageFunc is the function used for formatting
	// messages to send to chat channels/threads.
	FormatMessageFunc func(entry Entry, args FormatMessageArgs) string
	// FormatMessageArgs represents the data for formatting log
	// message.
	FormatMessageArgs struct {
		Service string
		Version string
		Prefix  string
	}
)

const (
	// FieldKey is the default key for saving fields
	// to the logger.
	FieldKey = "fields"
	// ErrorKey is the default key for saving errors
	// to the logger.
	ErrorKey = "error"
)

var (
	// DefaultReportFn is the default report function when
	// none is passed to the constructor.
	DefaultReportFn = func(e Entry) bool {
		return true
	}
	// DefaultFormatMessageFn is the default function used for sending
	// messages. It prints a formatted message from the log entry to
	// a user-friendly message.
	DefaultFormatMessageFn = func(entry Entry, args FormatMessageArgs) string {
		buf := bytes.Buffer{}

		// Write version from the latest build.
		buf.WriteString(fmt.Sprintf("%v %s", args.Service, emoji.ChartIncreasing))
		if args.Version != "" {
			buf.WriteString(fmt.Sprintf("v%s\n", strings.ReplaceAll(args.Version, "v", "")))
		}

		// Write intro text.
		app := strings.Title(strings.ToLower(args.Prefix)) //nolint
		buf.WriteString(fmt.Sprintf("\U0001FAE0 Error detected in %s, please see the information below for more details.\n\n", app))

		// Write log
		buf.WriteString(fmt.Sprintf("%v Level: %s\n", emoji.RightArrow, entry.Level))
		buf.WriteString(fmt.Sprintf("%v Time: %s\n", emoji.RightArrow, entry.Time.String()))
		if entry.Message != "" {
			buf.WriteString(fmt.Sprintf("%v Message: %s\n", emoji.RightArrow, entry.Message))
		}

		// Print out the Entries error.
		e := entry.Error()
		if e != nil {
			buf.WriteString(fmt.Sprintf("%v Code: %s\n", emoji.RightArrow, e.Code))
			buf.WriteString(fmt.Sprintf("%v Message: %s\n", emoji.RightArrow, e.Message))
			buf.WriteString(fmt.Sprintf("%v Operation: %s\n", emoji.RightArrow, e.Operation))
			buf.WriteString(fmt.Sprintf("%v Error: %s\n", emoji.RightArrow, e.Err))
			buf.WriteString(fmt.Sprintf("%v Fileline: %s\n\n", emoji.RightArrow, e.FileLine()))
		}

		// Print out associated data.
		if len(entry.Data) > 0 {
			buf.WriteString("Log entries:\n")
			for k, v := range entry.Data {
				if k == ErrorKey {
					continue
				}
				buf.WriteString(fmt.Sprintf("%s: %v\n", k, v))
			}
		}

		return buf.String()
	}
)

// ToLogrusEntry transforms an Entry to logrus.Entry
func (e Entry) ToLogrusEntry() logrus.Entry {
	return logrus.Entry(e)
}

// IsHTTP returns true if the fields contain HTTP
// key and value pairs.
func (e Entry) IsHTTP() bool {
	_, status := e.Data["status_code"]
	_, ip := e.Data["client_ip"]
	_, url := e.Data["request_url"]
	if status && ip && url {
		return true
	}
	return false
}

// Fields obtains the entries Fields if it exists,
// otherwise it returns nil.
func (e Entry) Fields() Fields {
	in, ok := e.Data[FieldKey]
	if !ok {
		return nil
	}
	fields, ok := in.(map[string]any)
	if !ok {
		return nil
	}
	return fields
}

// HasError determines if an error is attached to
// the entry.
func (e Entry) HasError() bool {
	err, ok := e.Data[ErrorKey]
	if !ok {
		return false
	}
	if err == nil {
		return false
	}
	return true
}

// Error returns a formatted error if one exists within the
// entry, otherwise returns nil.
func (e Entry) Error() *errors.Error {
	err, ok := e.Data[ErrorKey]
	if !ok {
		return nil
	}
	return errors.ToError(err)
}
