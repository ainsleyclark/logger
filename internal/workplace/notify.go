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
	"github.com/ainsleyclark/logger/types"
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
		Token         string
		Thread        string
		Service       string
		Version       string
		Prefix        string
		FormatMessage types.FormatMessageFunc
	}
)

// Fire will be called when some logging function is
// called with current hook. It will format log
// entry to string and write it to
// appropriate writer
func (hook *Hook) Fire(entry *logrus.Entry) error {
	go hook.process(types.Entry(*entry)) // This is already check for nil pointer.
	return nil
}

func (hook *Hook) process(entry types.Entry) {
	// Setup args for formatting the message.
	var (
		message = "" //nolint
		args    = types.FormatMessageArgs{
			Service: hook.options.Service,
			Version: hook.options.Version,
			Prefix:  hook.options.Prefix,
		}
	)

	// Use the default format message if none is attached,
	// otherwise call the function that is assigned.
	if hook.options.FormatMessage == nil {
		message = FormatMessage(entry, args)
	} else {
		message = hook.options.FormatMessage(entry, args)
	}

	// Use the Workplace client to send a message via the bot.
	err := hook.wp.Notify(workplace.Transmission{
		Thread:  hook.options.Thread,
		Message: message,
	})
	if err != nil {
		log.Println(err.Error()) // We can't use the standard logger as it may cause a loop.
	}
}

// FormatMessage prints a formatted message from the log entry to
// a user friendly message.
func FormatMessage(entry types.Entry, args types.FormatMessageArgs) string {
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
