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

package workplace

import (
	"bytes"
	"fmt"
	"github.com/ainsleyclark/errors"
	"github.com/ainsleyclark/mogrus"
	"github.com/ainsleyclark/workplace"
	"github.com/enescakir/emoji"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
)

// NewHook creates a new Workplace hook.
// Returns an error if the client could not be created.
func NewHook(opts Options) (*Hook, error) {
	wp, err := workplace.New(workplace.Config{Token: opts.Token})
	if err != nil {
		return nil, err
	}
	return &Hook{
		wp:        wp,
		options:   opts,
		LogLevels: logrus.AllLevels,
	}, nil
}

type (
	// Hook represents the workplace hook notifier
	// for log entries.
	Hook struct {
		wp        workplace.Notifier
		options   Options
		LogLevels []logrus.Level
	}
	// Options defines the configuration needed to fire logs
	// via the Workplace API.
	Options struct {
		Token   string
		Thread  string
		Service string
		Version string
		Prefix  string
	}
	FormatMessageFunc func(entry logrus.Entry)
)

// Fire will be called when some logging function is
// called with current hook. It will format log
// entry to string and write it to
// appropriate writer
func (hook *Hook) Fire(entry *logrus.Entry) error {
	go hook.process(mogrus.ToEntry(entry))
	return nil
}

func (hook *Hook) process(entry mogrus.Entry) {
	// Bail if the error is nil.
	if entry.Error == nil {
		return
	}

	// Bail if the error code is not anything but INTERNAL,
	// we don't want to notify users of invalid or pesky
	// log entries.
	if entry.Error.Code != errors.INTERNAL {
		return
	}

	// Use the Workplace client to send a message via the bot.
	err := hook.wp.Notify(workplace.Transmission{
		Thread:  hook.options.Thread,
		Message: hook.formatMessage(entry),
	})
	if err != nil {
		log.Println(err.Error()) // We can't use the standard logger as it may cause a loop.
	}
}

// formatMessage prints a formatted message from the log entry to
// a user friendly message.
func (hook *Hook) formatMessage(entry mogrus.Entry) string {
	buf := bytes.Buffer{}

	// Write version from the latest build.
	buf.WriteString(fmt.Sprintf("%v %s", hook.options.Service, emoji.ChartIncreasing))
	if hook.options.Version != "" {
		buf.WriteString(fmt.Sprintf("v%s\n", strings.ReplaceAll(hook.options.Version, "v", "")))
	}

	// Write intro text.
	app := strings.Title(strings.ToLower(hook.options.Prefix)) //nolint
	buf.WriteString(fmt.Sprintf("\U0001FAE0 Error detected in %s, please see the information below for more details.\n\n", app))

	// Write log
	buf.WriteString(fmt.Sprintf("%v Level: %s\n", emoji.RightArrow, entry.Level))
	buf.WriteString(fmt.Sprintf("%v Time: %s\n", emoji.RightArrow, entry.Time.String()))
	if entry.Message != "" {
		buf.WriteString(fmt.Sprintf("%v Message: %s\n", emoji.RightArrow, entry.Message))
	}

	// Print out the Entries error.
	buf.WriteString(fmt.Sprintf("%v Code: %s\n", emoji.RightArrow, entry.Error.Code)) // TODO: Handle nil pointer, don't print if the err is nil.
	buf.WriteString(fmt.Sprintf("%v Message: %s\n", emoji.RightArrow, entry.Error.Message))
	buf.WriteString(fmt.Sprintf("%v Operation: %s\n", emoji.RightArrow, entry.Error.Operation))
	buf.WriteString(fmt.Sprintf("%v Error: %s\n", emoji.RightArrow, entry.Error.Err))
	buf.WriteString(fmt.Sprintf("%v Fileline: %s\n\n", emoji.RightArrow, entry.Error.FileLine))

	// Print out associated data.
	if len(entry.Data) > 0 {
		buf.WriteString("Log entries:\n")
		for k, v := range entry.Data {
			if k == logrus.ErrorKey {
				continue
			}
			buf.WriteString(fmt.Sprintf("%s: %v\n", k, v))
		}
	}

	return buf.String()
}

// Levels Define on which log levels this hook would
// trigger.
func (hook *Hook) Levels() []logrus.Level {
	return hook.LogLevels
}
